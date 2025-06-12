package aws

import (
	"context"
	"os"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awsCredentials "github.com/aws/aws-sdk-go-v2/credentials"
	awsSes "github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/dpe27/es-krake/pkg/log"
)

type SesRepositoryInterface interface {
	SendEmail(
		ctx context.Context,
		from string,
		to string,
		subject string,
		contentText string,
		contentHtml string,
	) (*awsSes.SendEmailOutput, error)
	SendBulkTemplatedEmailWithContext(
		ctx context.Context,
		templateName string,
		defaultTemplateData string,
		from string,
		destinations []types.BulkEmailDestination,
	) (*awsSes.SendBulkTemplatedEmailOutput, error)

	CreateTemplateMail(
		ctx context.Context,
		templateName string,
		htmlBody string,
		htmlText string,
		subject string,
	) (*awsSes.CreateTemplateOutput, error)
	UpdateTemplateMail(
		ctx context.Context,
		templateName string,
		htmlBody string,
		htmlText string,
		subject string,
	) (*awsSes.UpdateTemplateOutput, error)
	DeleteTemplateMail(
		ctx context.Context,
		templateName string,
	) (*awsSes.DeleteTemplateOutput, error)
}

type SesRepository struct {
	service *awsSes.Client
}

// NewSesRepository comment
// en: create an instance of SesRepository
// en: with aws session and aws ses service
func NewSesRepository() (SesRepositoryInterface, error) {
	var err error
	var cfg aws.Config

	if os.Getenv("SES_AUTOCONFIG") == "true" {
		cfg, err = config.LoadDefaultConfig(
			context.Background(),
			config.WithRegion(os.Getenv("SES_REGION")),
			config.WithLogger(awsLogger{logger: log.With("service", "aws-ses-api")}),
			config.WithClientLogMode(aws.LogRequest|aws.LogResponseWithBody|aws.LogRetries),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(
			context.Background(),
			config.WithRegion(os.Getenv("SES_REGION")),
			config.WithCredentialsProvider(awsCredentials.NewStaticCredentialsProvider(
				os.Getenv("SES_CREDENTIALS_ID"),
				os.Getenv("SES_CREDENTIALS_SECRET"),
				os.Getenv("SES_CREDENTIALS_TOKEN"),
			)),
			config.WithLogger(awsLogger{logger: log.With("service", "aws-ses-api")}),
			config.WithClientLogMode(aws.LogRequest|aws.LogResponseWithBody|aws.LogRetries),
		)
	}
	if err != nil {
		return nil, err
	}

	client := awsSes.NewFromConfig(cfg)

	return &SesRepository{
		service: client,
	}, nil
}

func (sr *SesRepository) getConfigurationSetName() *string {
	name := os.Getenv("SES_CONFIGURATION_SET")
	if name == "" {
		return nil
	}
	return aws.String(name)
}

// SendEmail comment
// en: send single email through ses
func (sr *SesRepository) SendEmail(
	ctx context.Context,
	from string,
	to string,
	subject string,
	contentText string,
	contentHtml string,
) (*awsSes.SendEmailOutput, error) {
	body := &types.Body{}
	if contentText != "" {
		body.Text = &types.Content{
			Charset: aws.String("UTF-8"),
			Data:    aws.String(contentText),
		}
	}

	if contentHtml != "" {
		body.Html = &types.Content{
			Charset: aws.String("UTF-8"),
			Data:    aws.String(contentHtml),
		}
	}

	return sr.service.SendEmail(ctx, &awsSes.SendEmailInput{
		ConfigurationSetName: sr.getConfigurationSetName(),
		Destination: &types.Destination{
			ToAddresses: []string{
				to,
			},
		},
		Message: &types.Message{
			Body: body,
			Subject: &types.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(from),
	})
}

// SendBulkTemplatedEmailWithContext comment
// en: send bulk email with defined template :templateName
func (sr *SesRepository) SendBulkTemplatedEmailWithContext(
	ctx context.Context,
	templateName string,
	defaultTemplateData string,
	from string,
	destinations []types.BulkEmailDestination,
) (*awsSes.SendBulkTemplatedEmailOutput, error) {
	return sr.service.SendBulkTemplatedEmail(ctx, &awsSes.SendBulkTemplatedEmailInput{
		ConfigurationSetName: sr.getConfigurationSetName(),
		Destinations:         destinations,
		DefaultTemplateData:  aws.String(defaultTemplateData),
		Source:               aws.String(from),
		Template:             aws.String(templateName),
	})
}

// CreateTemplateMail comment
// en: create template mail in ses
func (sr *SesRepository) CreateTemplateMail(
	ctx context.Context,
	templateName string,
	htmlBody string,
	htmlText string,
	subject string,
) (*awsSes.CreateTemplateOutput, error) {
	return sr.service.CreateTemplate(ctx, &awsSes.CreateTemplateInput{
		Template: &types.Template{
			TemplateName: aws.String(templateName),
			HtmlPart:     aws.String(htmlBody),
			TextPart:     aws.String(htmlText),
			SubjectPart:  aws.String(subject),
		},
	})
}

// UpdateTemplateMail comment
// en: update template mail in ses
func (sr *SesRepository) UpdateTemplateMail(
	ctx context.Context,
	templateName string,
	htmlBody string,
	htmlText string,
	subject string,
) (*awsSes.UpdateTemplateOutput, error) {
	return sr.service.UpdateTemplate(ctx, &awsSes.UpdateTemplateInput{
		Template: &types.Template{
			TemplateName: aws.String(templateName),
			HtmlPart:     aws.String(htmlBody),
			TextPart:     aws.String(htmlText),
			SubjectPart:  aws.String(subject),
		},
	})
}

// DeleteTemplateMail comment
// en: delete template mail in ses
func (sr *SesRepository) DeleteTemplateMail(
	ctx context.Context,
	templateName string,
) (*awsSes.DeleteTemplateOutput, error) {
	return sr.service.DeleteTemplate(ctx, &awsSes.DeleteTemplateInput{
		TemplateName: aws.String(templateName),
	})
}
