package common

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const DEFAULT_SYNQ_EMAIL_SUBJECT = "SYNQ Email Notification"

type Email struct {
	ApiKey   string `json:"api_key"`
	Request  EmailRequest
	Response []EmailResponse
}

type EmailRequest struct {
	Sender        map[string]string `json:"sender"`
	Recipients    map[string]string `json:"recipients"`
	Subject       string            `json:"subject"`
	EmailBody     string            `json:"body"`
	EmailHtmlBody string            `json:"html_body"`
}

type EmailResponse struct {
	StatusCode int
	Body       string
	Headers    map[string][]string
}

func prepareEmailList(list map[string]string) (emails []*mail.Email) {
	for address, name := range list {
		emails = append(emails, mail.NewEmail(name, address))
	}
	return emails
}

func NewEmail() Email {
	return Email{
		Request:  EmailRequest{},
		Response: []EmailResponse{},
	}
}

func (e Email) Validate() (errors []error) {
	if e.ApiKey == "" {
		errors = append(errors, fmt.Errorf("Missing SendGrid API Key."))
	}

	if len(e.Request.Sender) < 1 {
		errors = append(errors, fmt.Errorf("Missing Sender value."))
	}

	if len(e.Request.Recipients) < 1 {
		errors = append(errors, fmt.Errorf("There must be one or more recipient addresses."))
	}

	if e.Request.EmailBody == "" {
		errors = append(errors, fmt.Errorf("The content value must be a string at least one character in length."))
	}

	return errors
}

func (e *Email) Send() []error {
	errors := e.Validate()
	if len(errors) > 0 {
		return errors
	}

	client := sendgrid.NewSendClient(e.ApiKey)
	from := prepareEmailList(e.Request.Sender)[0]
	recipients := prepareEmailList(e.Request.Recipients)
	subject := e.Request.Subject
	if subject == "" {
		subject = DEFAULT_SYNQ_EMAIL_SUBJECT
	}
	plainTextContent := e.Request.EmailBody
	htmlContent := e.Request.EmailHtmlBody
	if htmlContent == "" {
		htmlContent = plainTextContent
	}

	for _, to := range recipients {
		message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
		response, err := client.Send(message)
		if err != nil {
			errors = append(errors, err)
		} else {
			e.Response = append(e.Response, EmailResponse{
				StatusCode: response.StatusCode,
				Body:       response.Body,
				Headers:    response.Headers,
			})
		}
	}

	return errors
}
