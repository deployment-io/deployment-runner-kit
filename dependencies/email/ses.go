package email

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	ses_types "github.com/aws/aws-sdk-go-v2/service/ses/types"
	"log"
	"os"
)

type Ses struct {
	client *ses.Client
}

func NewSes(region string) (*Ses, error) {
	awsKeyID := os.Getenv("AWS_ACCESS_KEY_INDIA_ID")
	awsAccessKey := os.Getenv("AWS_SECRET_ACCESS_INDIA_KEY")
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsKeyID, awsAccessKey, "")))
	if err != nil {
		return nil, err
	}

	sesClient := ses.NewFromConfig(cfg, func(options *ses.Options) {
		options.Region = region
	})

	s := &Ses{client: sesClient}
	return s, nil
}

func convertAddressToStringArray(addresses []*Address) []string {
	var emails []string
	for _, address := range addresses {
		emails = append(emails, address.Email)
	}
	return emails
}

func (s Ses) SendEmail(from *Address, tos []*Address, bccs []*Address, subject *string,
	htmlContent *string, _ *string) error {

	destination := &ses_types.Destination{
		BccAddresses: convertAddressToStringArray(bccs),
		CcAddresses:  nil,
		ToAddresses:  convertAddressToStringArray(tos),
	}
	//"John Doe" <johndoe@example.com>
	source := from.Email
	if len(from.Name) > 0 {
		source = fmt.Sprintf("%s <%s>", from.Name, from.Email)
	}
	sendEmailInput := &ses.SendEmailInput{
		Destination: destination,
		Message: &ses_types.Message{
			Body: &ses_types.Body{
				Html: &ses_types.Content{
					Data:    htmlContent,
					Charset: nil, //we might need this in the future
				},
			},
			Subject: &ses_types.Content{
				Data:    subject,
				Charset: nil, //we might need this in the future
			},
		},
		Source:               aws.String(source),
		ConfigurationSetName: nil,
		ReplyToAddresses:     nil,
		ReturnPath:           nil,
		ReturnPathArn:        nil,
		SourceArn:            nil,
		Tags:                 nil,
	}
	sendEmailOutput, err := s.client.SendEmail(context.TODO(), sendEmailInput)
	if err != nil {
		return err
	}

	//TODO - remove later. for debugging for now.
	log.Println("ses message id: ", aws.ToString(sendEmailOutput.MessageId))

	return nil
}
