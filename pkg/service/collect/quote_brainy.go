package collect

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/gaaon/quote-collector/pkg/model"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type QuoteBrainyService struct {
	httpClient *http.Client
	baseUrl    string
	apiUrl     string
}

func NewQuoteBrainyService(httpClient *http.Client) *QuoteBrainyService {
	return &QuoteBrainyService{
		httpClient: httpClient,
		baseUrl:    "https://www.brainyquote.com/",
		apiUrl:     "https://www.brainyquote.com/api/inf",
	}
}

type QuotesContentResponse struct {
	Content string `json:"content"`
	Count   int    `json:"qCount"`
}

type QuotesContentRequest struct {
	Typ        string `json:"typ"`
	Lang       string `json:"langc"`
	V          string `json:"v"`
	Ab         string `json:"ab"`
	Pagination int    `json:"pg"`
	Id         string `json:"id"`
	Vid        string `json:"vid"`
	Fdd        string `json:"fdd"`
	M          int    `json:"m"`
}

func NewQuotesContentReq(vid string, pid string, pg int) *QuotesContentRequest {
	return &QuotesContentRequest{
		Typ:        "author",
		Lang:       "en",
		V:          "9.0.2:3290921",
		Ab:         "a",
		Pagination: pg,
		Id:         pid,
		Vid:        vid,
		Fdd:        "d",
		M:          0,
	}
}

func (service *QuoteBrainyService) findVidAndPersonId(link string) (
	vid string, pid string, err error) {

	req, err := http.NewRequest("GET", service.baseUrl+link, nil)
	if err != nil {
		return
	}

	res, err := service.httpClient.Do(req)
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

func (service *QuoteBrainyService) readQuotesFromReader(reader io.Reader) (quotes []model.Quote, err error) {
	docs, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return
	}

	quoteContents := docs.Find("div.qll-bg a.b-qt").Map(func(n int, s *goquery.Selection) string {
		return s.Text()
	})

	for _, content := range quoteContents {
		quotes = append(quotes, model.Quote{
			Content: content,
		})
	}

	return
}

func (service *QuoteBrainyService) findQuotesWithPagination(vid string, pid string, pg int) ([]model.Quote, error) {
	reqBody := NewQuotesContentReq(vid, pid, pg)

	reqBodyStr, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		service.apiUrl, bytes.NewReader(reqBodyStr))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := service.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New("status code is not ok")
	}

	bodyRaw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resBody QuotesContentResponse
	if err = json.Unmarshal(bodyRaw, &resBody); err != nil {
		return nil, err
	}

	return service.readQuotesFromReader(strings.NewReader(resBody.Content))
}

func (service *QuoteBrainyService) FindAllQuotesByLink(link string) (quotes []model.Quote, lastPagi int, err error) {
	vid, pid, err := service.findVidAndPersonId(link)
	if err != nil {
		return
	}

	var i int
	for i = 1; i < 100; i++ {
		var partialQuotes []model.Quote
		if partialQuotes, err = service.findQuotesWithPagination(vid, pid, i); err != nil {
			return
		}

		if len(partialQuotes) == 0 {
			break
		}

		quotes = append(quotes, partialQuotes...)
	}

	lastPagi = i - 1

	return
}
