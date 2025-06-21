package service

import (
	"bytes"
	"fmt"
	netMail "net/mail"
	textTemplate "text/template"

	"github.com/dpe27/es-krake/pkg/utils"
)

type SendEmailType int

const (
	SendEmailTypeHtml      SendEmailType = 1
	SendEmailTypeText      SendEmailType = 2
	SendEmailTypeMultipart SendEmailType = 3
	MaxLengthSender        int           = 18

	MaxReceiverPerSent = 50
)

type SendEmailParams struct {
	sendEmailType SendEmailType
	from          string
	to            string
	subject       string
	bodyHtml      string
	bodyText      string
}

func NewSendEmailParams(
	sendEmailType SendEmailType,
	from netMail.Address,
	to string,
	subject string,
	bodyHtml string,
	bodyText string,
) (*SendEmailParams, error) {
	if from.Address == "" {
		return nil, fmt.Errorf("invalid mailFrom: %#v", from)
	}
	from.Name = utils.SliceUTF8(from.Name, MaxLengthSender, "...")

	if to == "" {
		return nil, fmt.Errorf("invalid mailTo: %v", to)
	}
	if subject == "" {
		return nil, fmt.Errorf("invalid mailSubject: %v", subject)
	}
	switch sendEmailType {
	case SendEmailTypeHtml:
		if bodyHtml == "" {
			return nil, fmt.Errorf("invalid mailBody (HTML) : %v", bodyHtml)
		}
		if bodyText != "" {
			return nil, fmt.Errorf("invalid mailBody (TEXT) : %v", bodyText)
		}
	case SendEmailTypeText:
		if bodyHtml != "" {
			return nil, fmt.Errorf("invalid mailBody (HTML) : %v", bodyHtml)
		}
		if bodyText == "" {
			return nil, fmt.Errorf("invalid mailBody (TEXT) : %v", bodyText)
		}
	case SendEmailTypeMultipart:
		if bodyHtml == "" {
			return nil, fmt.Errorf("invalid mailBody (HTML) : %v", bodyHtml)
		}
		if bodyText == "" {
			return nil, fmt.Errorf("invalid mailBody (TEXT) : %v", bodyText)
		}
	default:
		return nil, fmt.Errorf("invalid sendEmailType: %v", sendEmailType)
	}
	return &SendEmailParams{
		sendEmailType: sendEmailType,
		from:          from.String(),
		to:            to,
		subject:       subject,
		bodyHtml:      bodyHtml,
		bodyText:      bodyText,
	}, nil
}

func NewSendEmailParamsFromTemplateFile(
	emailType SendEmailType,
	from netMail.Address,
	to string,
	subjectTemplate string,
	subjectValues interface{},
	bodyTemplateFileWithoutExtension string,
	bodyHtmlValues interface{},
	bodyTextValues interface{},
) (*SendEmailParams, error) {
	subject, err := renderTextMailTemplateFromString(
		subjectTemplate,
		subjectValues,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to render mail subject (template = %#v, data = %#v) %w",
			subjectTemplate,
			subjectValues,
			err,
		)
	}

	bodyHtml := ""
	if emailType == SendEmailTypeHtml || emailType == SendEmailTypeMultipart {
		mailTemplatePath := GetPathFileTemplate(bodyTemplateFileWithoutExtension + ".html")
		body, err := ReadMailTemplate(mailTemplatePath, bodyHtmlValues)
		if err != nil {
			return nil, fmt.Errorf("failed to read mail template '%v': %w", mailTemplatePath, err)
		}
		bodyHtml = body
	}
	bodyText, err := func() (string, error) {
		if emailType == SendEmailTypeText || emailType == SendEmailTypeMultipart {
			mailTemplatePath := GetPathFileTemplate(
				bodyTemplateFileWithoutExtension + ".txt",
			)
			body, err := renderTextMailTemplateFromPath(
				mailTemplatePath,
				bodyTextValues,
			)
			if err != nil {
				return "", fmt.Errorf(
					"failed to read mail template '%v': %w",
					mailTemplatePath,
					err,
				)
			}

			return body, nil
		}
		return "", nil
	}()
	if err != nil {
		return nil, err
	}

	return NewSendEmailParams(
		emailType,
		from,
		to,
		subject,
		bodyHtml,
		bodyText,
	)
}

func (p *SendEmailParams) SendEmailType() SendEmailType {
	return p.sendEmailType
}

func (p *SendEmailParams) From() string {
	return p.from
}

func (p *SendEmailParams) To() string {
	return p.to
}

func (p *SendEmailParams) Subject() string {
	return p.subject
}

func (p *SendEmailParams) BodyHtml() string {
	return p.bodyHtml
}

func (p *SendEmailParams) BodyText() string {
	return p.bodyText
}

func renderTextMailTemplateFromPath(
	templatePath string,
	data interface{},
) (string, error) {
	tmpl, err := textTemplate.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return "", err
	}

	return body.String(), nil
}

func renderTextMailTemplateFromString(
	template string,
	data interface{},
) (string, error) {
	tmpl, err := textTemplate.New("").Parse(template)
	if err != nil {
		return "", err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return "", err
	}

	return body.String(), nil
}

type SendBulkEmailParam struct {
	to                  string
	replaceTemplateData string
}

type TemplateParam struct {
	TemplateData string
	TemplateName string
}
