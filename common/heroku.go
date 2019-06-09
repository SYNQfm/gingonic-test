package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const HEROKU_BASE_URL = "https://api.heroku.com/"

type HerokuResponse struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

func UpdateHerokuVar(authKey, appName string, config interface{}) error {
	herokuUrl := HEROKU_BASE_URL + "apps/" + appName + "/config-vars"

	data, _ := json.Marshal(config)
	body := bytes.NewBuffer(data)

	req, err := makeRequest("PATCH", herokuUrl, authKey, body)
	if err != nil {
		return err
	}

	err = handleRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func makeRequest(method, url, token string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url, body)
	if err != nil {
		log.Println("could not create request: ", err.Error())
		return req, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/vnd.heroku+json; version=3")

	return req, nil
}

func handleRequest(req *http.Request) error {
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println("could not make http request: ", err.Error())
		return err
	}
	return parseResponse(resp)
}

func parseResponse(resp *http.Response) error {
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		responseAsBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("could not read resp body", err.Error())
			return err
		}

		hResp := HerokuResponse{}
		err = json.Unmarshal(responseAsBytes, &hResp)
		if err != nil {
			return err
		}

		return fmt.Errorf("%d: %s", resp.StatusCode, hResp.Message)
	}

	return nil
}
