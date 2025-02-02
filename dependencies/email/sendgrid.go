package email

import (
	"errors"
	"github.com/sendgrid/sendgrid-go"

	"strconv"

	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Sendgrid struct {
	apiKey string
}

func (s Sendgrid) SendEmail(from *Address, tos []*Address, bccs []*Address, subject *string,
	htmlContent *string, _ *string) error {
	f := mail.NewEmail(from.Name, from.Email)
	t := mail.NewEmail(tos[0].Name, tos[0].Email)
	plainTextContent := htmlContent
	message := mail.NewSingleEmail(f, *subject, t, *plainTextContent, *htmlContent)
	if len(tos) > 1 {
		for _, to := range tos {
			message.Personalizations[0].AddTos(mail.NewEmail(to.Name, to.Email))
		}
	}
	if bccs != nil && len(bccs) > 0 {
		for _, bcc := range bccs {
			message.Personalizations[0].AddBCCs(mail.NewEmail(bcc.Name, bcc.Email))
		}
	}
	client := sendgrid.NewSendClient(s.apiKey)
	res, err := client.Send(message)
	if err != nil {
		return err
	}
	if res.StatusCode > 300 {
		return errors.New("sendgrid returned " + strconv.Itoa(res.StatusCode))
	}
	return nil
}
