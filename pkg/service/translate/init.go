package translate

import "net/http"

var client *http.Client

func init() {
	client = &http.Client{}
}
