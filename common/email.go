package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

const ERIS_STAGING_URL = "https://4oiv2vbc7l.execute-api.eu-central-1.amazonaws.com/stage/eris"

type Email struct {
	EmailUrl string
	Request  EmailRequest
	Response EmailResponse
}

type EmailRequest struct {
	Service   string                 `json:"service"`
	EmailType string                 `json:"email_type"` // ERROR, INFO, or if template, use template ID
	Message   string                 `json:"message"`    // required for ERROR or INFO
	Body      map[string]interface{} `json:"body"`       // required except for ERROR and INFO
}

type EmailResponse struct {
	StatusCode int                 `json:"statusCode"`
	Message    string              `json:"message"`
	Errors     map[string][]string `json:"error"`
}

func NewEmail(emailUrl ...string) Email {
	// set eris staging as default email url
	var serviceUrl string
	if len(emailUrl) > 0 {
		serviceUrl = emailUrl[0]
	} else {
		serviceUrl = GetOsEnv("ERIS_URL", false, ERIS_STAGING_URL)
	}
	return Email{
		EmailUrl: serviceUrl,
		Request:  EmailRequest{},
		Response: EmailResponse{},
	}
}

func (e Email) Validate() error {
	if reflect.DeepEqual(e.Request, EmailRequest{}) {
		return fmt.Errorf("[Validate] Cannot send empty request!")
	}

	if e.Request.Service == "" {
		return fmt.Errorf("[Validate] Missing microservice name.")
	}

	if e.Request.EmailType == "" {
		return fmt.Errorf("[Validate] Missing email type.")
	}

	emailType := e.Request.EmailType
	if emailType == "ERROR" || emailType == "INFO" {
		if e.Request.Message == "" {
			return fmt.Errorf("[Validate] Request Message cannot be empty for ERROR or INFO email type.")
		}
	} else {
		if len(e.Request.Body) == 0 {
			return fmt.Errorf("[Validate] Request Body cannot be empty for %s email type.", emailType)
		}
	}

	return nil
}

func (e *Email) Send() error {
	err := e.Validate()
	if err != nil {
		return err
	}

	b, _ := json.Marshal(e.Request)
	body := bytes.NewBuffer(b)

	fmt.Printf("[INFO][Send] Calling email service %s with values: %s\n", e.EmailUrl, string(b))
	req, err := http.NewRequest("POST", e.EmailUrl, body)
	if err != nil {
		return fmt.Errorf("[Send] Error creating request %s: %s", string(b), err.Error())
	}
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("[Send] Error with request %s: %s", string(b), err.Error())
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("[Send] Error reading response for request %s: %s", string(b), err.Error())
	}

	err = json.Unmarshal(respBody, &e.Response)
	if err != nil {
		return fmt.Errorf("[Send] Error unmarshalling response %s: %s", string(respBody), err.Error())
	}

	return nil
}

func (e *Email) SetMsgAndSend(message string) error {
	message = strings.Replace(message, "\n", "<br>", -1)
	e.Request.Message = message
	return e.Send()
}
