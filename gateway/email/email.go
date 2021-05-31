package email

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"

	"go.uber.org/zap"
)

type EmailGateway interface {
	Send(to string, from string, payload string, subject string) (bool, error)
}

type emailGateway struct {
	SESSession *ses.SES
}

func New() EmailGateway {
	// Create an AWS session
	AWSRegion := os.Getenv("AWS_REGION")
	AWSSecretID := os.Getenv("AWS_SECRET_ID")
	AWSSecret := os.Getenv("AWS_SECRET")

	if len(AWSRegion) < 1 || len(AWSSecretID) < 1 || len(AWSSecret) < 1 {
		zap.L().Fatal("Please set required environment variables",
			zap.String("variables", "AWS_REGION, AWS_SECRET_ID, and AWS_SECRET"))
	}

	s, err := session.NewSession(&aws.Config{
		Region:      aws.String(AWSRegion),
		Credentials: credentials.NewStaticCredentials(AWSSecretID, AWSSecret, ""),
	})

	if err != nil {
		zap.L().Fatal("error creating AWS session", zap.Error(err))
	}

	// Create an SES session.
	ses := ses.New(s)

	return &emailGateway{SESSession: ses}
}

func (eg *emailGateway) Send(to string, from string, payload string, subject string) (bool, error) {
	sendEmail(eg.SESSession, to, from, payload, subject)
	return true, nil
}

func sendEmail(sess *ses.SES, to string, from string, payload string, subject string) (bool, error) {

	// Assemble the email
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(to),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String("<div>" + payload + "</div>"),
				},
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(payload),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(from),
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	result, err := sess.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				zap.L().Warn(ses.ErrCodeMessageRejected, zap.Error(aerr))
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				zap.L().Warn(ses.ErrCodeMailFromDomainNotVerifiedException, zap.Error(aerr))
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				zap.L().Warn(ses.ErrCodeConfigurationSetDoesNotExistException, zap.Error(aerr))
			default:
				zap.Error(aerr)
			}
		} else {
			zap.Error(aerr)
		}

		return false, err
	}

	zap.L().Info("email sent to address: "+to,
		zap.String("result", result.String()))
	return true, nil
}
