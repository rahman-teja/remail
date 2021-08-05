package remail

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type awsSes struct {
	region  string
	config  *aws.Config
	charSet string
	sender  string
}

func NewAWSSes(region, charSet, sender string, config *aws.Config) Remail {
	return &awsSes{
		region:  region,
		config:  config,
		charSet: charSet,
		sender:  sender,
	}
}

func (a awsSes) Send(ctx context.Context, messages Message) (err error) {
	if messages.From == "" {
		messages.From = a.sender
	}

	sess, err := session.NewSession(a.config)
	if err != nil {
		return err
	}

	// Create an SES session.
	svc := ses.New(sess)

	var body *ses.Body
	if messages.Body.ContentType == ContentTypeHTML {
		body = &ses.Body{
			Html: &ses.Content{
				Charset: aws.String(a.charSet),
				Data:    aws.String(mustBuildBody(messages.Body)),
			},
		}
	} else if messages.Body.ContentType == ContentTypePlaintext {
		body = &ses.Body{
			Text: &ses.Content{
				Charset: aws.String(a.charSet),
				Data:    aws.String(mustBuildBody(messages.Body)),
			},
		}
	}

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: aws.StringSlice(buildRecepient(messages.Cc)),
			ToAddresses: aws.StringSlice(buildRecepient(messages.To)),
		},
		Message: &ses.Message{
			Body: body,
			Subject: &ses.Content{
				Charset: aws.String(a.charSet),
				Data:    aws.String(messages.Subject),
			},
		},
		Source: aws.String(messages.From),
	}

	// Attempt to send the email.
	_, err = svc.SendEmailWithContext(ctx, input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				err = ErrCodeMessageRejected
				return
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				err = ErrCodeMailFromDomainNotVerifiedException
				return
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				err = ErrCodeConfigurationSetDoesNotExistException
				return
			default:
				err = aerr
				return
			}
		}
	}

	return
}
