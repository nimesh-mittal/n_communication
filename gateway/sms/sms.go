package sms

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"

	"go.uber.org/zap"
)

type SMSGateway interface {
	Send(to string, from string, payload string) (bool, error)
}

type smsGateway struct {
}

var twilioURL string
var twilioAccountSID string
var twilioAuthToken string

func New() SMSGateway {
	twilioAccountSID = os.Getenv("TWILIO_ACCOUNT_SID")
	twilioAuthToken = os.Getenv("TWILIO_AUTH_TOKEN")
	twilioURL = os.Getenv("TWILIO_URL")

	if len(twilioURL) < 1 || len(twilioAuthToken) < 1 || len(twilioAccountSID) < 1 {
		zap.L().Fatal("Please set required environment variables",
			zap.String("variables", "TWILIO_URL, TWILIO_AUTH_TOKEN, and TWILIO_ACCOUNT_SID"))
	}

	return &smsGateway{}
}

func (eg *smsGateway) Send(to string, from string, payload string) (bool, error) {
	return sendSMS(to, from, payload)
}

func sendSMS(to string, from string, payload string) (bool, error) {
	msgData := url.Values{}
	msgData.Set("To", to)
	msgData.Set("From", from)
	msgData.Set("Body", payload)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", twilioURL, &msgDataReader)
	req.SetBasicAuth(twilioAccountSID, twilioAuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)

	if err != nil {
		return false, err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err != nil {
			return false, err
		}

		zap.L().Info("sms sent successfully", zap.Any("sid", data["sid"]))
		return true, nil
	}

	zap.L().Info("error sending sms")
	return false, errors.New("received http status " + resp.Status)
}
