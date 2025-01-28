package email

import (
	"errors"
	"fmt"
	"github.com/deployment-io/deployment-runner-kit/enums/email_providers"
)

func GetProvider(p email_providers.Provider, opts ...*Options) (SendEmailR, error) {
	switch p {
	case email_providers.Sendgrid:
		if len(opts) == 0 {
			return nil, fmt.Errorf("api key missing for Sendgrid provider")
		}
		return Sendgrid{
			apiKey: *opts[0].ApiKey,
		}, nil
	case email_providers.Ses:
		if len(opts) == 0 {
			return nil, fmt.Errorf("region missing for Ses provider")
		}
		ses, err := NewSes(*opts[0].AwsRegion)
		if err != nil {
			return nil, err
		}
		return ses, nil
	case email_providers.Smtp:
		if len(opts) == 0 {
			return nil, fmt.Errorf("smtp host missing for Smtp provider")
		}
		smtp, err := NewSmtp(*opts[0].SmtpHost, *opts[0].SmtpPort, *opts[0].SmtpUsername, *opts[0].SmtpPassword)
		if err != nil {
			return nil, err
		}
		return smtp, nil
	}
	return nil, errors.New("invalid provider")
}
