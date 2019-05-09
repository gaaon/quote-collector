package collect

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
)

func GetKoreanNameFromEnglish(name string) (koreanName string, err error){
	urlStr := "https://google.co.kr/search?ie=UTF-8&q=" + url.QueryEscape(name)
	var (
		req *http.Request
		res *http.Response
		doc *goquery.Document
	)
	if req, err = http.NewRequest("GET", urlStr, nil); err != nil {
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.140 Safari/537.36 Edge/17.17134")

	if res, err = client.Do(req); err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode == 429 {
		fmt.Printf("%+v\n", res.Header)
	}
	if doc, err = goquery.NewDocumentFromReader(res.Body); err != nil {
		return
	}

	println(res.StatusCode)

	doc.Find(".kno-fb-ctx.gsmt").Each(func(idx int, el *goquery.Selection) {
		koreanName = el.Text()
	})

	return
}
