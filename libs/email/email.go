package email

import (
	"errors"
	"fmt"
	"net/smtp"
	"strings"
	"sync"
)

type EmailType int

const (
	html EmailType = 1
	text EmailType = 2
)

type EmailValidate interface {
	Send(to string, body string, subject string, emailType EmailType) error
}

type EmailClient struct {
	addr   string
	from   string
	name   string
	secret string
	pool   sync.Pool
}

type unencryptedAuth struct {
	smtp.Auth
}

func (a unencryptedAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	s := *server
	s.TLS = true
	return a.Auth.Start(&s)
}

func New(addr, from, name, secret string) *EmailClient {
	e := &EmailClient{
		addr:   addr,
		from:   from,
		name:   name,
		secret: secret,
	}
	e.pool.New = func() interface{} {
		client, err := smtp.Dial(addr)
		if err != nil {
			return nil
		}
		return client
	}
	return e
}

func (email *EmailClient) Send(to string, body string, subject string, emailType EmailType) error {
	hp := strings.Split(email.addr, ":")
	auth := unencryptedAuth{
		smtp.PlainAuth(
			"",
			email.from,
			email.secret,
			hp[0],
		),
	}
	c, ok := email.pool.Get().(*smtp.Client)
	if !ok {
		return errors.New("send email failed")
	}
	defer email.pool.Put(c)
	if err := c.Auth(auth); err != nil {
		return err
	}
	if err := c.Mail(email.from); err != nil {
		return err
	}
	if err := c.Rcpt(to); err != nil {
		return err
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	data := email.generateMessage(to, body, subject, emailType)
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return nil
	}
	return c.Quit()
}

func (e *EmailClient) generateMessage(to string, body string, subject string, emailType EmailType) []byte {
	header := make(map[string]string)
	header["From"] = e.name + "<" + e.from + ">"
	header["To"] = to
	header["Subject"] = subject

	if emailType == 1 {
		header["Content-Type"] = "text/html; charset=UTF-8"
	} else if emailType == 2 {
		header["Content-Type"] = "text/plain; charset=UTF-8"
	}

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s:%s\r\n", k, v)
	}
	message += "\r\n" + body

	return []byte(message)
}
