package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

type Smtp struct {
	host     string
	port     string
	username string
	password string
}

func NewSmtp(smtpHost string, port string, username string, password string) (*Smtp, error) {
	s := &Smtp{
		host:     smtpHost,
		port:     port,
		username: username,
		password: password,
	}
	return s, nil
}

func (s Smtp) SendEmail(from *Address, tos []*Address, _ []*Address, subject *string, htmlContent *string, plainTextContent *string) error {
	// Ensure at least one content part exists; otherwise, return an error
	if plainTextContent == nil && htmlContent == nil {
		return fmt.Errorf("both plainTextContent and htmlContent are nil; cannot send an empty email")
	}
	if subject == nil {
		return fmt.Errorf("subject is nil; cannot send an email without a subject")
	}
	if from == nil {
		return fmt.Errorf("from is nil; cannot send an email without a from address")
	}
	if len(tos) == 0 {
		return fmt.Errorf("tos is empty; cannot send an email without a to address")
	}
	// SMTP server details
	smtpHost := s.host
	smtpPort := s.port
	username := s.username
	password := s.password

	// Sender and recipient details
	f := from.Email
	recipientList := make([]string, len(tos))
	for i, recipient := range tos {
		recipientList[i] = recipient.Email
	}

	// MIME Headers
	fromHeader := fmt.Sprintf("From: %s\n", f)
	// Join the email addresses with commas
	toHeader := fmt.Sprintf("To: %s\n", strings.Join(recipientList, ", "))
	subjectHeader := fmt.Sprintf("Subject: %s\n", *subject)

	// Boundary for the multipart/alternative content
	boundary := "boundary-example-123456"

	// MIME Headers for multipart content
	contentTypeHeader := fmt.Sprintf("MIME-Version: 1.0\nContent-Type: multipart/alternative; boundary=%s\n\n", boundary)

	// Initialize email body parts
	var plainTextPart string
	var htmlPart string

	// Add plain text part only if it is not nil
	if plainTextContent != nil {
		plainTextPart = fmt.Sprintf(
			"--%s\nContent-Type: text/plain; charset=\"utf-8\"\n\n%s\n\n",
			boundary, *plainTextContent,
		)
	}

	// Add HTML part only if it is not nil
	if htmlContent != nil {
		htmlPart = fmt.Sprintf(
			"--%s\nContent-Type: text/html; charset=\"utf-8\"\n\n%s\n\n",
			boundary, *htmlContent,
		)
	}

	// Closing boundary
	closingBoundary := fmt.Sprintf("--%s--", boundary)

	// Combine all parts into the final message
	message := []byte(fromHeader + toHeader + subjectHeader + contentTypeHeader + plainTextPart + htmlPart + closingBoundary)

	// Authentication for SMTP
	auth := smtp.PlainAuth("", username, password, smtpHost)

	// Send the email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, f, recipientList, message)
	if err != nil {
		return err
	}

	return nil
}

var _ SendEmailR = Smtp{}
