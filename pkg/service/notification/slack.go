package notification

import (
	"errors"
	"net/http"
	"os"
	"strings"
)

type SlackService struct {
	httpClient  *http.Client
	hookUrl string
}

func NewSlackService(httpClient *http.Client) (*SlackService, error){
	var slackWebhookUrl string

	if slackWebhookUrl = os.Getenv("SLACK_WEBHOOK_URL"); slackWebhookUrl == "" {
		return nil, errors.New("slack webhook url cannot be empty")
	}
	return &SlackService{
		httpClient: httpClient,
		hookUrl: slackWebhookUrl,
	}, nil
}

func (service *SlackService) SendNotiToDevice(message string) error {
	payload := "{\"text\": \"" + message + "\"}"

	req, err := http.NewRequest("POST", service.hookUrl, strings.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := service.httpClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New("slack webhook api not 200 status")
	}

	return nil
}