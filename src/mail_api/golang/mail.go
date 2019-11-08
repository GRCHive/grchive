package mail

import "errors"

var ErrorMailDeliveryFailed error = errors.New("Mail delivery mailed.")

type MailAPIProvider string

const (
	SendgridProvider MailAPIProvider = "sendgrid"
)

type Email struct {
	Name  string
	Email string
}

type MailPayload struct {
	From    Email
	To      Email
	Subject string
	Message string
}

type MailProviderInterface interface {
	Initialize(args ...interface{})
	SendMail(payload MailPayload) error
}

var MailProvider MailProviderInterface = nil

func InitializeMailAPI(provider MailAPIProvider, args ...interface{}) {
	switch provider {
	case SendgridProvider:
		MailProvider = new(SendgridInterface)
	}

	if MailProvider == nil {
		panic("Unknown Mail API Provider")
		return
	}

	MailProvider.Initialize(args...)
}
