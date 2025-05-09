package mailer

import "context"

// Email is email notification
type Email struct {
	Sender     string
	Recipients []string
	CC         []string
	BCC        []string
	Subject    string
	Attachment []byte
	Message    string
}

type Mailer interface {
	Send(context.Context, Email) error
}
