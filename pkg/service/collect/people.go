package collect

import (
	"github.com/gaaon/quote-collector/pkg/model"
	"github.com/gaaon/quote-collector/pkg/repository"
	"io"
	"net/http"
	"sort"
	"sync"
)

/**
	find people list started with nameChar from a page
 */
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

			if partialPeopleList, err = repository.FindPeopleListByReaderWithAnchor(bodyReader); err != nil {
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

	index := 0
	for i, person := range peopleList {
		if person.FullName[0] < 'A' || person.FullName[0] > 'Z' {
			index = i
			break
		}
	}

	return peopleList[0:index], nil
}

func FindPeopleListFromSnapshot() (peopleList []model.Person, err error) {
	if repository.IsExistPeopleListSnapshot() {
		if peopleList, err = FindPeopleListFromAToZ(); err != nil {
			return
		}

		if err = repository.SavePeopleListIntoSnapshot(peopleList); err != nil {
			return
		}

		println("Read from server")

		return
	} else {
		if peopleList, err = repository.FindPeopleListFromSnapshot(); err != nil {
			return
		}

		println("Read from local")

		return
	}
}