package collect

import (
	"net/http"
	"time"
)

var client *http.Client

func init() {
	client = &http.Client{
		Timeout: time.Second * 10,
	}
}
