package notification

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type PushoverService struct {
	client  *http.Client
	userKey string
	token   string
}

func NewPushoverService(httpClient *http.Client) (*PushoverService, error) {
	var (
		userKey string
		token   string
	)

	if userKey = os.Getenv("PUSHOVER_USER_KEY"); userKey == "" {
		return nil, errors.New("pushover user key cannot be empty")
	}

	if token = os.Getenv("PUSHOVER_TOKEN"); token == "" {
		return nil, errors.New("pushover token cannot be empty")
	}

	return &PushoverService{
		httpClient,
		userKey,
		token,
	}, nil
}

func (service *PushoverService) SendNotiToDevice(message string) error {
	values := url.Values{}
	values.Add("token", service.token)
	values.Add("user", service.userKey)
	values.Add("message", message)
	values.Add("title", "quote-collector")

	req, err := http.NewRequest(
		"POST",
		"https://api.pushover.net/1/messages.json",
		strings.NewReader(values.Encode()))

	if err != nil {
		return err
	}

	res, err := service.client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New("error caused while sending noti message")
	}

	return nil
}
