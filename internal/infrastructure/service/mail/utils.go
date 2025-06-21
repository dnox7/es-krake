package service

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"os"
	"regexp"
)

func GetSystemEmailFrom() string {
	return os.Getenv("SYSTEM_EMAIL_FROM")
}

// return file path of template mail
func GetPathFileTemplate(filePath string) string {
	localPath := os.Getenv("DW_TEMPLATES_PATH")
	return localPath + filePath
}

// return body mail
func ReadMailTemplate(templatePath string, data interface{}) (string, error) {
	var body bytes.Buffer
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}
	err = t.Execute(&body, data)
	if err != nil {
		return "", err
	}
	return body.String(), err
}

func ReadMailTemplateWithData(templateMaildata *string, data interface{}) (string, error) {
	var body bytes.Buffer

	bodyMail := ""
	if templateMaildata != nil {
		bodyMail = *templateMaildata
		regexpNewLine := regexp.MustCompile(`\r?\n`)
		bodyMail = regexpNewLine.ReplaceAllString(bodyMail, "<br>")
	}

	t, err := template.New("mail").Parse(bodyMail)
	if err != nil {
		return "", err
	}

	err = t.Execute(&body, data)
	if err != nil {
		return "", err
	}
	return body.String(), err
}

func GetCommentReportEmailTo() string {
	return os.Getenv("COMMENT_REPORT_EMAIL_TO")
}

func GetAllAlertEmailTo() string {
	return os.Getenv("ALL_ALERT_EMAIL_TO")
}

// Builds a `from` value using a name and email
// https://docs.aws.amazon.com/ses/latest/DeveloperGuide/send-email-raw.html#send-email-mime-encoding-headers
func BuildMailFrom(name, email string) string {
	base64Name := base64.StdEncoding.EncodeToString([]byte(name))
	encodedName := "=?UTF-8?B?" + base64Name + "?="
	return fmt.Sprintf("%v <%v>", encodedName, email)
}
