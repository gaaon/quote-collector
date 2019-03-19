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
	"time"
)

var client *http.Client

func init() {
	client = &http.Client{
		Timeout: time.Second * 10,
	}
}

type CompositeName struct {
	TitleName string
	TextName  string
}

type CompositeNames []CompositeName

func (c CompositeNames) Len() int {
	return len(c)
}

func (c CompositeNames) Swap(i, j int) {
	c[i], c[j] = c[j], c[i];
}

func (c CompositeNames) Less(i, j int) bool {
	if c[i].TitleName == c[j].TitleName {
		return c[i].TextName < c[j].TextName
	}

	return c[i].TitleName < c[j].TitleName
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

func GetPeopleListByReader(bodyReader io.ReadCloser) (peopleList []CompositeName, err error) {
	defer bodyReader.Close()

	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		return
	}

	doc.Find("h3 .mw-headline").Parent().Next().Each(func(_ int, peopleUl *goquery.Selection) {
		peopleUl.Find("li a").Each(func(_ int, nameLink *goquery.Selection) {
			name := strings.TrimSpace(nameLink.Text())
			peopleList = append(peopleList, CompositeName {
				name,
				nameLink.Text()})
		})
	})

	return
}

func GetPeopleListByReaderWithAnchor(bodyReader io.ReadCloser) (peopleList []CompositeName, err error) {
	defer bodyReader.Close()

	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		return
	}

	doc.Find("h3 .mw-headline").Parent().Next().Each(func(_ int, peopleUl *goquery.Selection) {
		peopleUl.Find("li a").Each(func(_ int, nameLink *goquery.Selection) {
			titleAttr, exists := nameLink.Attr("title")
			if exists {
				name := strings.TrimSpace(titleAttr)
				peopleList = append(peopleList, CompositeName {
					name,
					nameLink.Text()})
			}
		})
	})

	return
}

var nameRanges = []string{
	"A", "B", "C", "D", "E-F", "G", "H", "I-J", "K", "L", "M", "N-O", "P", "Q-R", "S", "T-V", "W-Z",
}
func GetPeopleListFromAToZ() (peopleList []CompositeName, err error) {
	var wg sync.WaitGroup
	var mutex sync.Mutex

	for _, nameRange := range nameRanges {
		wg.Add(1)
		go func(_nameRange string) {
			var (
				bodyReader io.ReadCloser
				partialPeopleList []CompositeName
			)

			if bodyReader, err = getPeopleListHtmlByName(_nameRange); err != nil {
				return
			}

			if partialPeopleList, err = GetPeopleListByReaderWithAnchor(bodyReader); err != nil {
				return
			}

			mutex.Lock()
			peopleList = append(peopleList, partialPeopleList...)
			mutex.Unlock()

			wg.Done()
		}(nameRange)
	}

	wg.Wait()

	sort.Sort(CompositeNames(peopleList))

	index := 0
	for i, name := range peopleList {
		if name.TitleName[0] < 'A' || name.TitleName[0] > 'Z' {
			index = i
			break
		}
	}
	return peopleList[0: index], nil
}

func saveIntoSnapshot(peopleList []CompositeName, fileVersion string) (err error) {
	f, err := os.Create("data/" + fileVersion + "/snapshot.txt")
	if err != nil {
		return
	}
	defer f.Close()

	for i, name := range peopleList {
		var newLineOrNot = ""
		if i != len(peopleList) - 1 {
			newLineOrNot = "\n"
		}

		_, _ = f.Write([]byte(name.TitleName + "\t" + name.TextName + newLineOrNot))
	}

	return
}

func GetPeopleListFromSnapshot(fileVersion string) (peopleList []CompositeName, err error) {
	f, err := os.Open("data/" + fileVersion + "/snapshot.txt")
	defer f.Close()

	if os.IsNotExist(err) {
		if peopleList, err = GetPeopleListFromAToZ(); err != nil {
			return
		}

		if err = saveIntoSnapshot(peopleList, fileVersion); err != nil {
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
			peopleList = append(peopleList, CompositeName{names[0], names[1]})
		}

		return
	}
}