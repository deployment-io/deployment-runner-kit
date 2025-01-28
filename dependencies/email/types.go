package email

type EmailI interface {
	SendEmailR
}

type SendEmailR interface {
	SendEmail(from *Address, tos []*Address, bccs []*Address, subject *string,
		htmlContent *string, plainTextContent *string) error
}

type Options struct {
	ApiKey       *string
	AwsRegion    *string
	SmtpHost     *string
	SmtpPort     *string
	SmtpUsername *string
	SmtpPassword *string
}

type Address struct {
	Name  string
	Email string
}
