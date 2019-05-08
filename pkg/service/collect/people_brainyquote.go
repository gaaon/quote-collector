package collect

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/gaaon/quote-collector/pkg/model"
	"github.com/gaaon/quote-collector/pkg/repository"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const brainyPeopleStartByUrl = "https://www.brainyquote.com/authors/"

var brainyHttpClient *http.Client

func init() {
	brainyHttpClient = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func FindPeopleListInBrainyByReader(bodyReader io.Reader) (peopleList []model.Person, err error) {
	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		return nil, err
	}

	doc.Find("div.bq_s > table > tbody> tr").Each(func (a int, s *goquery.Selection) {
		onClickValue, _ := s.Attr("onclick")

		link := strings.Trim(
			strings.TrimRight(strings.Split(onClickValue, "=")[1], ";"),
			"'")

		fullName := s.Find("td > a").Text()

		peopleList = append(peopleList, model.Person{
			FullName: fullName,
			Link: link,
		})
	})

	return
}

func FindPeopleListInBrainyStartsWith(startChar string, pagination int) (
	[]model.Person, error) {

		brainyPeopleListUrl := brainyPeopleStartByUrl + startChar
		if pagination != 1 {
			brainyPeopleListUrl += strconv.Itoa(pagination)
		}

		req, err := http.NewRequest(
			"GET",
			brainyPeopleListUrl,
			nil)
		if err != nil {
			return nil, err
		}

		res, err := brainyHttpClient.Do(req)
		if err != nil {
			return nil, err
		}

		if res.StatusCode == 301 || res.StatusCode == 404 { // page not exists
			return nil, errors.New("page not exists")
		}

		defer res.Body.Close()

		return FindPeopleListInBrainyByReader(res.Body)
}

func FindPeopleListInBrainyFromSnapshot() (peopleList []model.Person, err error){
	if repository.IsExistPeopleListSnapshot() {
		println("Read brainy people list from local")
		return repository.FindPeopleListFromSnapshot()
	}

	println("Read brainy people list from server")

	var wg sync.WaitGroup
	var mutex sync.Mutex
	for i := 'a'; i <= 'z'; i++ {
		wg.Add(1)

		go func(startChar int32) {
			for j := 1; ; j++ {
				partialPeopleList, err := FindPeopleListInBrainyStartsWith(string(startChar), j)
				if err != nil {
					if err.Error() == "page not exists" {
						break
					}
				}

				mutex.Lock()
				peopleList = append(peopleList, partialPeopleList...)
				mutex.Unlock()
			}

			wg.Done()
		}(i)
	}

	wg.Wait()

	sort.Sort(model.PeopleSorts(peopleList))

	err = repository.SavePeopleListIntoSnapshot(peopleList)

	return
}