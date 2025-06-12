package aws

import (
	"context"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awsCredentials "github.com/aws/aws-sdk-go-v2/credentials"
	awsSes "github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	appCfg "github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/pkg/log"
)

type SesService interface {
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

type Ses struct {
	cli                  *awsSes.Client
	configurationSetName string
}

func NewSesService(cfg *appCfg.Config) (*Ses, error) {
	awsCfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(cfg.SES.Region),
		config.WithCredentialsProvider(awsCredentials.NewStaticCredentialsProvider(
			cfg.SES.CredentialsID,
			cfg.SES.CredentialsSecret,
			cfg.SES.CredentialsToken,
		)),
		config.WithLogger(awsLogger{logger: log.With("service", "aws-ses-api")}),
		config.WithClientLogMode(aws.LogRequest|aws.LogResponseWithBody|aws.LogRetries),
	)
	if err != nil {
		return nil, err
	}

	client := awsSes.NewFromConfig(awsCfg)

	return &Ses{
		cli: client,
	}, nil
}

// SendEmail comment
// en: send single email through ses
func (sr *Ses) SendEmail(
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

	return sr.cli.SendEmail(ctx, &awsSes.SendEmailInput{
		ConfigurationSetName: &sr.configurationSetName,
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
func (sr *Ses) SendBulkTemplatedEmailWithContext(
	ctx context.Context,
	templateName string,
	defaultTemplateData string,
	from string,
	destinations []types.BulkEmailDestination,
) (*awsSes.SendBulkTemplatedEmailOutput, error) {
	return sr.cli.SendBulkTemplatedEmail(ctx, &awsSes.SendBulkTemplatedEmailInput{
		ConfigurationSetName: &sr.configurationSetName,
		Destinations:         destinations,
		DefaultTemplateData:  aws.String(defaultTemplateData),
		Source:               aws.String(from),
		Template:             aws.String(templateName),
	})
}

// CreateTemplateMail comment
// en: create template mail in ses
func (sr *Ses) CreateTemplateMail(
	ctx context.Context,
	templateName string,
	htmlBody string,
	htmlText string,
	subject string,
) (*awsSes.CreateTemplateOutput, error) {
	return sr.cli.CreateTemplate(ctx, &awsSes.CreateTemplateInput{
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
func (sr *Ses) UpdateTemplateMail(
	ctx context.Context,
	templateName string,
	htmlBody string,
	htmlText string,
	subject string,
) (*awsSes.UpdateTemplateOutput, error) {
	return sr.cli.UpdateTemplate(ctx, &awsSes.UpdateTemplateInput{
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
func (sr *Ses) DeleteTemplateMail(
	ctx context.Context,
	templateName string,
) (*awsSes.DeleteTemplateOutput, error) {
	return sr.cli.DeleteTemplate(ctx, &awsSes.DeleteTemplateInput{
		TemplateName: aws.String(templateName),
	})
}
