package collect

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gaaon/quote-collector/pkg/model"
	"io/ioutil"
	"net/http"
	"strings"
)

const brainyQuoteBaseUrl = "https://www.brainyquote.com/"
const brainyQuoteApiUrl = "https://www.brainyquote.com/api/inf"

type QuotesContentResponse struct {
	Content string `json:"content"`
	Count int `json:"qCount"`
}

type QuotesContentRequest struct {
	Typ string `json:"typ"`
	Lang string `json:"langc"`
	V string `json:"v"`
	Ab string `json:"ab"`
	Pagination int `json:"pg"`
	Id string `json:"id"`
	Vid string `json:"vid"`
	Fdd string `json:"fdd"`
	M int `json:"m"`
}

func NewQuotesContentReq(vid string, pid string, pg int) *QuotesContentRequest {
	return &QuotesContentRequest{
		Typ: "author",
		Lang: "en",
		V: "9.0.2:3290921",
		Ab: "a",
		Pagination: pg,
		Id: pid,
		Vid: vid,
		Fdd: "d",
		M: 0,
	}
}

func FindVidAndPersonIdInBrainy(path string) (
	vid string, pid string, err error) {

	req, err := http.NewRequest("GET", brainyQuoteBaseUrl + path, nil)
	if err != nil {
		return
	}

	res, err := client.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()

	rawBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	lines := strings.Split(string(rawBody), "\n")

	for _, line := range lines {
		if strings.Contains(line, "VID=") {
			vidWrapper := strings.Split(line, "=")[1]
			vidWrapper = strings.TrimRight(vidWrapper, ";")
			vid = strings.Trim(vidWrapper, "'")
		} else if strings.Contains(line, "ctTarg[\"bq_aId\"]") {
			pidWrapper := strings.TrimSpace(strings.Split(line, "=")[1])
			pidWrapper = strings.TrimRight(pidWrapper, ";")
			pid = strings.Trim(pidWrapper, "\"")
		}
	}
	return
}

func FindQuotesInBrainy(vid string, pid string, pg int) ([]model.Quote, error) {
	reqBody := NewQuotesContentReq(vid, pid, pg)

	reqBodyStr, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		brainyQuoteApiUrl, bytes.NewReader(reqBodyStr))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	bodyRaw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	println(string(bodyRaw))
	var resBody QuotesContentResponse
	if err = json.Unmarshal(bodyRaw, &resBody); err != nil {
		return nil, err
	}

	fmt.Printf("%+v\n", resBody)

	return nil, nil
}
//func FindQuotesInBrainyByLinkPath(path string) ([]model.Quote, error) {
//	req, err := http.NewRequest("GET", brainyQuoteBaseUrl + path, nil)
//	if err != nil {
//		return nil, err
//	}
//
//	res, err := client.Do(req)
//	if err != nil {
//		return nil, err
//	}
//
//	defer res.Body.Close()
//
//	docs, err := goquery.NewDocumentFromReader(res.Body)
//	if err != nil {
//		return nil, err
//	}
//
//
//}
