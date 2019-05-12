package collect

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
)

type NameTranslateService struct {
	httpClient *http.Client
}

func NewNameTranslateService(httpClient *http.Client) *NameTranslateService {
	return &NameTranslateService{
		httpClient: httpClient,
	}
}

func (service *NameTranslateService) TranslateFullNameToKorean(fullName string) (koreanName string, err error){
	urlStr := "https://google.co.kr/search?ie=UTF-8&q=" + url.QueryEscape(fullName)
	var (
		req *http.Request
		res *http.Response
		doc *goquery.Document
	)

	if req, err = http.NewRequest("GET", urlStr, nil); err != nil {
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.140 Safari/537.36 Edge/17.17134")

	if res, err = service.httpClient.Do(req); err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode == 429 {
		fmt.Printf("429 status header : %+v\n", res.Header)

		return "", errors.New("too many request status code from server")
	}
	if doc, err = goquery.NewDocumentFromReader(res.Body); err != nil {
		return
	}

	doc.Find(".kno-fb-ctx.gsmt").Each(func(idx int, el *goquery.Selection) {
		koreanName = el.Text()
	})

	return
}
