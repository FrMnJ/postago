package email

import (
	"bytes"
	"html/template"
	"log"
	"path"

	"github.com/FrMnJ/postago/src/config"
	"github.com/wneessen/go-mail"
)

type GmailEmailServiceAdapter struct {
	FromAddress  string
	Passwd       string
	ServerDomain string
	ServerPort   int
}

func NewGmailEmailServiceAdapter() *GmailEmailServiceAdapter {
	return &GmailEmailServiceAdapter{
		FromAddress:  config.AppConfig.MailConfig.Account,
		Passwd:       config.AppConfig.MailConfig.Passwd,
		ServerDomain: config.AppConfig.MailConfig.SmtpHost,
		ServerPort:   config.AppConfig.MailConfig.SmtpPort,
	}
}

func (es *GmailEmailServiceAdapter) SendEmail(
	toEmail, subject, templateName string,
	data map[string]any,
) error {
	m, err := es.NewMessage(toEmail)
	if err != nil {
		return err
	}
	m.Subject(subject)

	baseTemplate, err := NewTemplate(templateName)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	if err := baseTemplate.Execute(&body, data); err != nil {
		return err
	}

	m.SetBodyString(mail.TypeTextHTML, body.String())

	c, err := es.NewClient()
	if err != nil {
		return err
	}
	return SendMessage(m, c)
}

func (es *GmailEmailServiceAdapter) NewMessage(toEmail string) (*mail.Msg, error) {
	m := mail.NewMsg()
	if err := m.From(es.FromAddress); err != nil {
		return nil, err
	}
	if err := m.To(toEmail); err != nil {
		return nil, err
	}
	return m, nil
}

func SendMessage(m *mail.Msg, c *mail.Client) error {
	if err := c.DialAndSend(m); err != nil {
		log.Fatalf("Failed to send mail: %v", err)
		return err
	}
	return nil
}

func (es *GmailEmailServiceAdapter) NewClient() (*mail.Client, error) {
	c, err := mail.NewClient(
		es.ServerDomain,
		mail.WithPort(es.ServerPort),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(es.FromAddress),
		mail.WithPassword(es.Passwd),
	)
	if err != nil {
		log.Fatalf("Failed to create mail client: %v", err)
		return nil, err
	}
	return c, nil
}

func NewTemplate(
	templateName string,
) (*template.Template, error) {
	t, err := template.ParseFiles(
		path.Join(
			config.GetBaseProjectPath(),
			"internal",
			"external",
			"email",
			"templates",
			"template.html",
		),
	)
	if err != nil {
		return nil, err
	}
	_, err = t.ParseFiles(
		path.Join(
			config.GetBaseProjectPath(),
			"src",
			"queue",
			"email",
			"templates",
			templateName,
		),
	)
	if err != nil {
		return nil, err
	}

	return t, nil
}
