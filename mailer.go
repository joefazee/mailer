package mailer

import (
	"bytes"
	"embed"
	"text/template"
	"time"

	"github.com/go-mail/mail/v2"
)

//go:embed "templates"
var templateFS embed.FS

type MailerConfig struct {
	Timeout      time.Duration
	Host         string
	Port         int
	Username     string
	Password     string
	Sender       string
	TemplatePath string
}

type Mailer struct {
	dailer *mail.Dialer
	config MailerConfig
	sender string
}

func New(config MailerConfig) Mailer {

	dailer := mail.NewDialer(config.Host, config.Port, config.Username, config.Password)
	dailer.Timeout = config.Timeout

	return Mailer{
		dailer: dailer,
		sender: config.Sender,
		config: config,
	}
}

func (m Mailer) Send(to, templateFile string, data interface{}) error {

	if m.config.TemplatePath == "" {
		m.config.TemplatePath = "templates/"
	}

	tmpl, err := template.New("email").ParseFS(templateFS, m.config.TemplatePath+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	msg := mail.NewMessage()
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject.String())
	msg.SetHeader("From", m.sender)
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	return m.dailer.DialAndSend(msg)
}
