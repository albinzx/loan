package mailer

import (
	"context"
	"fmt"
	"log"
	"net/smtp"

	"github.com/albinzx/loan/pkg/config"
)

// SMTPServer is smtp server configuration
type SMTPServer struct {
	Address string
	Sender  string
	Auth    smtp.Auth
}

type SMTPService struct {
	config config.Config
}

// New return new SMTP service
func New(config config.Config) *SMTPService {
	return &SMTPService{
		config: config,
	}
}

func (es *SMTPService) Send(ctx context.Context, email Email) error {
	user := es.config.GetString("smtp.user")
	pass := es.config.GetString("smtp.password")
	host := es.config.GetString("smtp.host")
	port := es.config.GetString("smtp.port")

	auth := smtp.PlainAuth("", user, pass, host)

	server := SMTPServer{
		Address: fmt.Sprintf("%s:%s", host, port),
		Auth:    auth,
	}

	if err := sendEmail(server, email); err != nil {
		log.Printf("error while sending email, %v", err)
		return err
	}

	return nil
}

func sendEmail(server SMTPServer, email Email) error {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + email.Subject + "\n"
	to := "To: " + email.Recipients[0] + "\n"
	msg := []byte(to + subject + mime + "\n" + email.Message)

	return smtp.SendMail(server.Address, server.Auth, server.Sender, email.Recipients, msg)
}
