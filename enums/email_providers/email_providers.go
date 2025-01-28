package email_providers

type Provider int

const (
	Sendgrid Provider = iota + 1
	Ses
	Smtp
)
