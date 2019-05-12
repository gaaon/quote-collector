package collect

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/gaaon/quote-collector/pkg/constant"
	"github.com/gaaon/quote-collector/pkg/model"
	"github.com/hashicorp/go-multierror"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type PeopleBrainyService struct {
	authorBaseUrl         string
	httpClient            *http.Client
	peopleSnapshotService *peopleSnapshotService
}

func NewPeopleBrainyService(peopleSnapshotService *peopleSnapshotService) (*PeopleBrainyService, error) {
	if peopleSnapshotService == nil {
		return nil, errors.New("peopleSnapshotService cannot be nil")
	}

	return &PeopleBrainyService{
		authorBaseUrl: "https://www.brainyquote.com/authors/",
		httpClient: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		peopleSnapshotService: peopleSnapshotService,
	}, nil
}

func (service *PeopleBrainyService) findPeopleListByBody(bodyReader io.Reader) (
	peopleList []model.Person, err error) {

	var doc *goquery.Document
	if doc, err = goquery.NewDocumentFromReader(bodyReader); err != nil {
		return
	}

	doc.Find("div.bq_s > table > tbody> tr").Each(func(a int, s *goquery.Selection) {
		onClickValue, _ := s.Attr("onclick")

		link := strings.Trim(
			strings.TrimRight(strings.Split(onClickValue, "=")[1], ";"),
			"'")

		fullName := s.Find("td > a").Text()

		peopleList = append(peopleList, model.Person{
			FullName: fullName,
			Link:     link,
			Source:   constant.SOURCE_BRAINY_QUOTE,
		})
	})

	return
}

func (service *PeopleBrainyService) findPeopleListStartsWith(
	nameFirstChar string, pagination int) (peopleList []model.Person, err error) {

	brainyPeopleListUrl := service.authorBaseUrl + nameFirstChar
	if pagination != 1 {
		brainyPeopleListUrl += strconv.Itoa(pagination)
	}

	var (
		req *http.Request
		res *http.Response
	)

	if req, err = http.NewRequest("GET", brainyPeopleListUrl, nil); err != nil {
		return
	}

	if res, err = service.httpClient.Do(req); err != nil {
		return
	}

	if res.StatusCode == 301 || res.StatusCode == 404 { // page not exists
		return
	}

	defer res.Body.Close()

	return service.findPeopleListByBody(res.Body)
}

func (service *PeopleBrainyService) FindPeopleListFromSnapshotOrRemote() (peopleList []model.Person, err error) {

	peopleSnapshotService := service.peopleSnapshotService

	if peopleSnapshotService.IsPeopleListExist() {
		log.Info("read brainy people list from local file")

		peopleList, err = peopleSnapshotService.FindPeopleList()
	} else {
		log.Info("read brainy people list from server")

		var wg sync.WaitGroup
		var mutex sync.Mutex
		for i := 'a'; i <= 'z'; i++ {
			wg.Add(1)

			go func(nameFirstChar int32) {
				defer wg.Done()

				for j := 1; ; j++ {
					partialPeopleList, err2 := service.findPeopleListStartsWith(string(nameFirstChar), j)
					if err2 != nil {
						mutex.Lock()
						err = multierror.Append(err, err2)
						mutex.Unlock()

						break
					}

					if len(partialPeopleList) == 0 {
						break
					}

					mutex.Lock()
					peopleList = append(peopleList, partialPeopleList...)
					mutex.Unlock()
				}
			}(i)
		}

		wg.Wait()
		if err != nil {
			return
		}

		sort.Sort(model.PeopleSorts(peopleList))

		if err = peopleSnapshotService.SavePeopleList(peopleList); err != nil {
			return
		}
	}

	return
}
