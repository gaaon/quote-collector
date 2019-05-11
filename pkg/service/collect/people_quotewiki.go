package collect

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gaaon/quote-collector/pkg/model"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
)

func FindPeopleListHtmlByName(nameChar string) (body io.ReadCloser, err error) {
	req, err := http.NewRequest(
		"GET",
		"https://en.wikiquote.org/wiki/List_of_people_by_name,_"+nameChar,
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

var nameRanges = []string{
	"A", "B", "C", "D", "E-F", "G", "H", "I-J", "K", "L", "M", "N-O", "P", "Q-R", "S", "T-V", "W-Z",
}

func FindPeopleListFromAToZ() (peopleList []model.Person, err error) {
	var wg sync.WaitGroup
	var mutex sync.Mutex

	for _, nameRange := range nameRanges {
		wg.Add(1)
		go func(_nameRange string) {
			var (
				bodyReader        io.ReadCloser
				partialPeopleList []model.Person
				err               error
			)

			if bodyReader, err = FindPeopleListHtmlByName(_nameRange); err != nil {
				return
			}
			defer bodyReader.Close()

			if partialPeopleList, err = FindPeopleListByReaderWithAnchor(bodyReader); err != nil {
				return
			}

			mutex.Lock()
			peopleList = append(peopleList, partialPeopleList...)
			mutex.Unlock()

			wg.Done()
		}(nameRange)
	}

	wg.Wait()

	sort.Sort(model.PeopleSorts(peopleList))

	return peopleList, nil
}

/**
	find title attribute in anchor if exists from a page
 */
func FindPeopleListByReaderWithAnchor(bodyReader io.ReadCloser) (peopleList []model.Person, err error) {
	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		return
	}

	doc.Find("h3 .mw-headline").Parent().Next().Each(func(_ int, peopleUl *goquery.Selection) {
		peopleUl.Find("li a").Each(func(_ int, nameLink *goquery.Selection) {
			titleAttr, exists := nameLink.Attr("title")
			linkUrl, _ := nameLink.Attr("href")

			if exists {
				name := strings.TrimSpace(titleAttr)
				peopleList = append(peopleList, model.Person{
					FullName:     name,
					ReversedName: nameLink.Text(),
					Link:         linkUrl,
				})
			}
		})
	})

	return
}

