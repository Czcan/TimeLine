package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

type EmailType int

const (
	html EmailType = 1
	text EmailType = 2
)

type EmailClient struct {
	addr   string
	from   string
	name   string
	secret string
}

func New(addr, from, name, secret string) (*EmailClient, error) {
	e := &EmailClient{
		addr:   addr,
		from:   from,
		name:   name,
		secret: secret,
	}
	return e, nil
}

func (client *EmailClient) Send(to string, body string, subject string, emailType EmailType) error {
	hp := strings.Split(client.addr, ":")
	auth := smtp.PlainAuth("", client.from, client.secret, hp[0])
	header := make(map[string]string)
	header["From"] = client.name + "<" + client.from + ">"
	header["To"] = to
	header["Subject"] = subject
	if emailType == 1 {
		header["Content-Type"] = "text/html; charset=UTF-8"
	} else if emailType == 2 {
		header["Content-Type"] = "text/plain; charset=UTF-8"
	}
	sendTo := strings.Split(to, ";")
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s:%s\r\n", k, v)
	}
	message += "\r\n" + body
	err := smtp.SendMail(client.addr, auth, client.from, sendTo, []byte(message))
	if err != nil {
		return err
	}
	return nil
}
