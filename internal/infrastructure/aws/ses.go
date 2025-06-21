package aws

import (
	"context"
	"crypto/rand"
	"errors"
	"math"
	"math/big"
	"strings"
	"time"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awsCredentials "github.com/aws/aws-sdk-go-v2/credentials"
	awsSes "github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	appCfg "github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/utils"
)

const (
	defaultRetryTime = 2
	errThrottle      = "Maximum sending rate exceeded"
)

type ResultAndIndex struct {
	result *awsSes.SendBulkTemplatedEmailOutput
	index  []int
}

type SesService interface {
	SendEmail(
		ctx context.Context,
		from string,
		to string,
		subject string,
		contentText string,
		contentHtml string,
	) error

	SendBulkTemplatedEmail(
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
	) error

	UpdateTemplateMail(
		ctx context.Context,
		templateName string,
		htmlBody string,
		htmlText string,
		subject string,
	) error

	DeleteTemplateMail(
		ctx context.Context,
		templateName string,
	) error
}

type Ses struct {
	configurationSetName string
	logger               *log.Logger
	cli                  *awsSes.Client
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
		cli:    client,
		logger: log.With("service", "ses-service"),
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
) error {
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

	resp, err := sr.cli.SendEmail(ctx, &awsSes.SendEmailInput{
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

	sr.log(ctx, to, subject, resp, err)
	return err
}

// SendBulkTemplatedEmail comment
// en: send bulk email with defined template :templateName
func (sr *Ses) SendBulkTemplatedEmail(
	ctx context.Context,
	templateName string,
	defaultTemplateData string,
	from string,
	destinations []types.BulkEmailDestination,
) (*awsSes.SendBulkTemplatedEmailOutput, error) {
	result := &awsSes.SendBulkTemplatedEmailOutput{
		Status: make([]types.BulkEmailDestinationStatus, len(destinations)),
	}
	currentLoop := 0
	var currentErr error
	currentIdx := utils.RangeN(len(destinations) - 1)

	var resultList []ResultAndIndex
	for currentLoop <= defaultRetryTime && len(destinations) > 0 {
		if currentLoop > 0 {
			time.Sleep(sr.exponentialDuration(ctx, currentLoop))
		}
		currentLoop++

		bulkResult, err := sr.sendBulkTemplatedEmailWithContext(
			ctx,
			templateName,
			defaultTemplateData,
			from,
			destinations,
		)
		if err != nil {
			if !sr.isThrottleError(err) {
				return bulkResult, err
			}
			// if send email has error, resend for all current destinations
			currentErr = err
			continue
		}

		// if send mail don't return error, we record result and resend mail only for destinations has throttle error
		resultList = append(resultList, ResultAndIndex{
			result: bulkResult,
			index:  currentIdx,
		})

		currentIdx, destinations = sr.getThrottleDestinationsIteration(bulkResult, destinations, currentIdx)
	}

	// if result list is empty, return current error
	if len(resultList) == 0 {
		return nil, currentErr
	}

	// has result, we iterate result list to get final result
	for _, resultAndIdx := range resultList {
		currentResult := resultAndIdx.result
		indexList := resultAndIdx.index
		if currentResult.Status != nil {
			for i := range currentResult.Status {
				result.Status[indexList[i]] = currentResult.Status[i]
			}
		}
	}

	return result, nil
}

// sendBulkTemplatedEmailWithContext comment
// en: send bulk email with defined template :templateName
func (sr *Ses) sendBulkTemplatedEmailWithContext(
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

func (sr *Ses) isThrottleError(err error) bool {
	var nsk *types.LimitExceededException
	return errors.As(err, &nsk)
}

func (sr *Ses) exponentialDuration(ctx context.Context, count int) time.Duration {
	timeDuration, err := rand.Int(rand.Reader, big.NewInt(500))
	if err != nil {
		sr.logger.Error(ctx, "error when random exponential duration", "error", err)
	}
	return time.Duration(int(math.Pow(2, float64(count))))*time.Second + time.Duration(timeDuration.Int64())*time.Millisecond
}

func (sr *Ses) getThrottleDestinationsIteration(
	bulkResponse *awsSes.SendBulkTemplatedEmailOutput,
	destinations []types.BulkEmailDestination,
	currentIdx []int,
) ([]int, []types.BulkEmailDestination) {
	var errMsg string
	var indexArr []int
	var destinationsArr []types.BulkEmailDestination
	for i, response := range bulkResponse.Status {
		err := response.Error
		if err != nil {
			errMsg = *err
		} else {
			errMsg = ""
		}

		if strings.Contains(errMsg, errThrottle) {
			destinationsArr = append(destinationsArr, destinations[i])
			indexArr = append(indexArr, currentIdx[i])
		}
	}

	return indexArr, destinationsArr
}

// CreateTemplateMail comment
// en: create template mail in ses
func (sr *Ses) CreateTemplateMail(
	ctx context.Context,
	templateName string,
	htmlBody string,
	htmlText string,
	subject string,
) error {
	_, err := sr.cli.CreateTemplate(ctx, &awsSes.CreateTemplateInput{
		Template: &types.Template{
			TemplateName: aws.String(templateName),
			HtmlPart:     aws.String(htmlBody),
			TextPart:     aws.String(htmlText),
			SubjectPart:  aws.String(subject),
		},
	})
	return err
}

// UpdateTemplateMail comment
// en: update template mail in ses
func (sr *Ses) UpdateTemplateMail(
	ctx context.Context,
	templateName string,
	htmlBody string,
	htmlText string,
	subject string,
) error {
	_, err := sr.cli.UpdateTemplate(ctx, &awsSes.UpdateTemplateInput{
		Template: &types.Template{
			TemplateName: aws.String(templateName),
			HtmlPart:     aws.String(htmlBody),
			TextPart:     aws.String(htmlText),
			SubjectPart:  aws.String(subject),
		},
	})
	return err
}

// DeleteTemplateMail comment
// en: delete template mail in ses
func (sr *Ses) DeleteTemplateMail(
	ctx context.Context,
	templateName string,
) error {
	_, err := sr.cli.DeleteTemplate(ctx, &awsSes.DeleteTemplateInput{
		TemplateName: aws.String(templateName),
	})
	return err
}

func (sr *Ses) log(
	ctx context.Context,
	to string,
	subject string,
	respObj *awsSes.SendEmailOutput,
	respErr error,
) {
	logger := sr.logger.With(
		"to", to,
		"subject", subject,
	)

	if respObj.MessageId != nil {
		logger = logger.With("messageId", *respObj.MessageId)
	}

	if respErr != nil {
		logger.Error(ctx, "send email failed", "error", respErr)
	} else {
		logger.Info(ctx, "send email success")
	}
}
