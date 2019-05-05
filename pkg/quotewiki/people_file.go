package quotewiki

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
)

type PeopleSorts []Person

func (c PeopleSorts) Len() int {
	return len(c)
}

func (c PeopleSorts) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c PeopleSorts) Less(i, j int) bool {
	if c[i].FullName == c[j].FullName {
		return c[i].ReversedName < c[j].ReversedName
	}

	return c[i].FullName < c[j].FullName
}

var nameRanges = []string{
	"A", "B", "C", "D", "E-F", "G", "H", "I-J", "K", "L", "M", "N-O", "P", "Q-R", "S", "T-V", "W-Z",
}

/**
	find people list started with nameChar from a page
 */
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

/**
	find title attribute in anchor if exists from a page
 */
func getPeopleListByReaderWithAnchor(bodyReader io.ReadCloser) (peopleList []Person, err error) {
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
				peopleList = append(peopleList, Person {
					FullName: name,
					ReversedName: nameLink.Text(),
					Link: linkUrl,
				})
			}
		})
	})

	return
}

func getPeopleListFromAToZ() (peopleList []Person, err error) {
	var wg sync.WaitGroup
	var mutex sync.Mutex

	for _, nameRange := range nameRanges {
		wg.Add(1)
		go func(_nameRange string) {
			var (
				bodyReader io.ReadCloser
				partialPeopleList []Person
				err error
			)

			if bodyReader, err = getPeopleListHtmlByName(_nameRange); err != nil {
				return
			}
			defer bodyReader.Close()

			if partialPeopleList, err = getPeopleListByReaderWithAnchor(bodyReader); err != nil {
				return
			}

			mutex.Lock()
			peopleList = append(peopleList, partialPeopleList...)
			mutex.Unlock()

			wg.Done()
		}(nameRange)
	}

	wg.Wait()

	sort.Sort(PeopleSorts(peopleList))

	index := 0
	for i, person := range peopleList {
		if person.FullName[0] < 'A' || person.FullName[0] > 'Z' {
			index = i
			break
		}
	}

	return peopleList[0: index], nil
}

const peopleListSnapshotLocation = "data/snapshot.txt"

func saveIntoWriter(peopleList []Person, writer io.Writer) (err error) {
	for i, name := range peopleList {
		var newLineOrNot = ""
		if i != len(peopleList) - 1 {
			newLineOrNot = "\n"
		}

		_, _ = writer.Write([]byte(
			name.FullName + "\t" +
				name.ReversedName + "\t" +
				name.Link + newLineOrNot))
	}

	return
}

func saveIntoSnapshot(peopleList []Person) (err error) {
	f, err := os.Create(peopleListSnapshotLocation)
	if err != nil {
		return
	}
	defer f.Close()

	return saveIntoWriter(peopleList, f)
}

func FindPeopleListFromSnapshot() (peopleList []Person, err error) {
	f, err := os.Open(peopleListSnapshotLocation)
	defer f.Close()

	if os.IsNotExist(err) {
		if peopleList, err = getPeopleListFromAToZ(); err != nil {
			return
		}

		if err = saveIntoSnapshot(peopleList); err != nil {
			return
		}

		println("Read from server")

		return
	} else {
		var content []byte
		if content, err = ioutil.ReadAll(f); err != nil {
			return
		}

		println("Read from local")

		splits := strings.Split(string(content), "\n")
		for _, split := range splits {
			names := strings.Split(split, "\t")
			peopleList = append(peopleList, Person{FullName: names[0], ReversedName: names[1], Link: names[2]})
		}

		return
	}
}