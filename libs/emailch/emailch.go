package email

import (
	"net/smtp"
	"strings"
	"time"

	"github.com/jordan-wright/email"
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
	p      *email.Pool
}

func New(addr, from, name, secret string) (*EmailClient, error) {
	e := &EmailClient{
		addr:   addr,
		from:   from,
		name:   name,
		secret: secret,
	}
	var err error
	h := strings.Split(addr, ":")
	e.p, err = email.NewPool(
		addr,
		5,
		smtp.PlainAuth("", from, secret, h[0]),
	)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (client *EmailClient) Send(to string, body string, subject string, emailType EmailType) error {
	e := email.NewEmail()
	e.From = client.name + "<" + client.from + ">"
	e.To = []string{to}
	e.Subject = subject
	if emailType == 1 {
		e.Text = []byte(body)
	} else if emailType == 2 {
		e.HTML = []byte(body)
	}
	return client.p.Send(e, time.Second*10)
}
