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

	var exists bool
	userKey, exists = os.LookupEnv("PUSHOVER_USER_KEY")
	if !exists {
		log.Fatal("no pushover user key")
	}

	token, exists = os.LookupEnv("PUSHOVER_TOKEN")
	if !exists {
		log.Fatal("no pushover token")
	}
}

func SendNotiToDevice(message string, title string) error {
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