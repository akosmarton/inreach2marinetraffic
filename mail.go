package main

import (
	"fmt"
	"net/smtp"
)

type SmtpClient struct {
	Host     string
	Port     string
	User     string
	Password string
}

type Mail struct {
	From    string
	To      string
	Subject string
	Body    string
}

func (c *SmtpClient) Send(m *Mail) error {
	auth := smtp.PlainAuth(
		"",
		c.User,
		c.Password,
		c.Host,
	)

	header := make(map[string]string)
	header["From"] = m.From
	header["To"] = m.To
	header["Subject"] = m.Subject
	header["Content-Type"] = "text/plain; charset=\"utf-8\""

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + m.Body

	return smtp.SendMail(
		fmt.Sprintf("%s:%s", c.Host, c.Port),
		auth,
		m.From,
		[]string{m.To},
		[]byte(message),
	)
}
