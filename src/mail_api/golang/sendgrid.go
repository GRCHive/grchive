package mail

import (
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"net/http"
)

const baseApiUrl string = "https://api.sendgrid.com"
const sendEmailEndpoint string = "/v3/mail/send"

type SendgridInterface struct {
	ApiKey string
}

func (s *SendgridInterface) Initialize(args ...interface{}) {
	// Expected args:
	// 	0 - API Key
	if len(args) != 1 {
		panic("Invalid SendGrid interface arguments.")
		return
	}
	s.ApiKey = args[0].(string)
}

func (s SendgridInterface) SendMail(payload MailPayload) error {
	from := mail.NewEmail(payload.From.Name, payload.From.Email)
	to := mail.NewEmail(payload.To.Name, payload.To.Email)
	message := mail.NewSingleEmail(from, payload.Subject, to, payload.Message, payload.Message)
	client := sendgrid.NewSendClient(s.ApiKey)
	resp, err := client.Send(message)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusAccepted {
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Body)
		fmt.Println(resp.Headers)
		return ErrorMailDeliveryFailed
	}

	return nil
}
