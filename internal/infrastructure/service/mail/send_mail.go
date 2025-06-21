package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"

	awsSdk "github.com/aws/aws-sdk-go-v2/aws"
	awsSes "github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/dpe27/es-krake/internal/infrastructure/aws"
)

type MailService interface {
	SendEmail(
		ctx context.Context,
		params *SendEmailParams,
	) error

	SendBulkEmail(
		ctx context.Context,
		templateParams TemplateParam,
		bulkEmailParams []SendBulkEmailParam,
		from string,
	) (*awsSes.SendBulkTemplatedEmailOutput, error)

	GetBodyAndSubjectTemplate(
		templatePath string,
		subject string,
	) (string, string, string, error)
}

type mailService struct {
	ses aws.SesService
}

func NewMailService(b aws.SesService) MailService {
	return &mailService{
		ses: b,
	}
}

// SendEmail implements MailService.
func (m *mailService) SendEmail(ctx context.Context, params *SendEmailParams) error {
	if params == nil {
		return errors.New("params is nil")
	}

	switch params.SendEmailType() {
	case SendEmailTypeHtml:
		return m.ses.SendEmail(ctx, params.From(), params.To(), params.Subject(), "", params.BodyHtml())
	case SendEmailTypeText:
		return m.ses.SendEmail(ctx, params.From(), params.To(), params.Subject(), params.BodyText(), "")
	case SendEmailTypeMultipart:
		return m.ses.SendEmail(ctx, params.From(), params.To(), params.Subject(), params.BodyText(), params.BodyHtml())
	default:
		return fmt.Errorf("invalid sendEmailType: %v", params.SendEmailType())
	}
}

// SendBulkEmail implements MailService.
func (m *mailService) SendBulkEmail(
	ctx context.Context,
	templateParams TemplateParam,
	bulkEmailParams []SendBulkEmailParam,
	from string,
) (*awsSes.SendBulkTemplatedEmailOutput, error) {
	batchResult := awsSes.SendBulkTemplatedEmailOutput{
		Status: []types.BulkEmailDestinationStatus{},
	}
	var destinations []types.BulkEmailDestination
	// loop through input, send 50 mail each time.
	for idx, emailParam := range bulkEmailParams {
		destinations = append(destinations, types.BulkEmailDestination{
			Destination: &types.Destination{
				ToAddresses: []string{emailParam.to},
			},
			ReplacementTemplateData: awsSdk.String(emailParam.replaceTemplateData),
		})
		if (idx+1)%MaxReceiverPerSent == 0 || idx == len(bulkEmailParams)-1 {
			result, err := m.ses.SendBulkTemplatedEmail(ctx, templateParams.TemplateName, templateParams.TemplateData, from, destinations)
			if err != nil {
				return nil, err
			}

			batchResult.Status = append(batchResult.Status, result.Status...)

			// reset destination
			destinations = []types.BulkEmailDestination{}

			// sleep one second between each send bulk email
			if idx < len(bulkEmailParams)-1 {
				time.Sleep(1 * time.Second)
			}
		}
	}

	return &batchResult, nil
}

// GetBodyAndSubjectTemplate implements MailService.
func (m *mailService) GetBodyAndSubjectTemplate(templatePath string, subject string) (string, string, string, error) {
	re := regexp.MustCompile(`{{\.\w+?}}`)

	transformHandlebarFormat := func(match string) string {
		return match[0:2] + match[3:]
	}

	subject = re.ReplaceAllStringFunc(subject, transformHandlebarFormat)

	bodyHtml, err := os.ReadFile(GetPathFileTemplate(templatePath + ".html"))
	if err != nil {
		return "", "", "", err
	}
	bodyHtmlHandlebar := re.ReplaceAllStringFunc(string(bodyHtml), transformHandlebarFormat)

	bodyText, err := os.ReadFile(GetPathFileTemplate(templatePath + ".html"))
	if err != nil {
		return "", "", "", err
	}
	bodyTextHandlebar := re.ReplaceAllStringFunc(string(bodyText), transformHandlebarFormat)

	return subject, bodyHtmlHandlebar, bodyTextHandlebar, nil
}
