package notification

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var client *http.Client
var userKey string
var token string

func init() {
	client = &http.Client{}

	userKey = os.Getenv("PUSHOVER_USER_KEY")

	token = os.Getenv("PUSHOVER_TOKEN")
}

func SendNotiToDevice(message string, title string) error {
	if userKey == "" {
		log.Fatal("user key cannot be empty")
	}

	if token == "" {
		log.Fatal("token cannot be empty")
	}

	values := url.Values{}
	values.Add("token", token)
	values.Add("user", userKey)
	values.Add("message", message)
	values.Add("title", title)

	req, err := http.NewRequest(
		"POST",
		"https://api.pushover.net/1/messages.json",
		strings.NewReader(values.Encode()))

	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New("error caused while sending noti message")
	}

	return nil
}