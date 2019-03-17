package quotewiki

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"strings"
	"time"
)

var client *http.Client

func init() {
	client = &http.Client{
		Timeout: time.Second * 10,
	}
}

func getPeopleListHtmlByName(nameChar string) (body io.ReadCloser, err error) {
	req, err := http.NewRequest(
		"GET",
		"https://en.wikiquote.org/wiki/List_of_people_by_name,_" + nameChar,
		nil)

	if err != nil {
		return
	}

	res, err := client.Do(req)
	if err != nil {
		return
	}

	return res.Body, nil
}

func GetPeopleListByName(bodyReader io.ReadCloser, nameChar string) (peopleList []string, err error) {
	defer bodyReader.Close()

	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		return
	}

	doc.Find("h3 .mw-headline").Parent().Next().Each(func(_ int, peopleUl *goquery.Selection) {
		peopleUl.Find("li a").Each(func(_ int, nameLink *goquery.Selection) {
			name := strings.TrimSpace(nameLink.Text())
			peopleList = append(peopleList, name)
		})
	})

	return
}