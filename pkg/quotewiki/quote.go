package quotewiki

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type Revision struct {
	XMLName xml.Name `xml:"revision"`
	Text string `xml:"text"`
}
type Page struct {
	XMLName xml.Name `xml:"page"`
	Title string `xml:"title"`
	Revision Revision `xml:"revision"`
}

type MediaWiki struct {
	XMLName xml.Name `xml:"mediawiki"`
	Pages []Page `xml:"page"`
}

var mediaWiki *MediaWiki

func getMediaWikiFromFile() (*MediaWiki, error) {
	if mediaWiki != nil {
		return mediaWiki, nil
	}

	xmlFile, err := os.Open("enwikiquote-latest-pages-articles.xml")
	if err != nil {
		return nil, err
	}

	defer xmlFile.Close()

	var xmlBody MediaWiki
	byteValue, _ := ioutil.ReadAll(xmlFile)
	err = xml.Unmarshal(byteValue, &xmlBody)
	if err != nil {
		return nil, err
	}

	mediaWiki = &xmlBody
	return mediaWiki, err
}

var _personPageMap map[string]*Page

func getTextFieldByTitleName(name string) (err error) {
	mediaWiki, err := getMediaWikiFromFile()
	if err != nil {
		return err
	}

	peopleList, err := FindPeopleListFromSnapshot()
	if err != nil {
		return err
	}

	personPageMap := make(map[string]*Page)

	for i, page := range mediaWiki.Pages {
		personPageMap[page.Title] = &mediaWiki.Pages[i]
	}

	cnt := 0
	quoteSub := 0

	peopleChan := make(chan Person, len(peopleList))
	for _, person := range peopleList {
		peopleChan <- person
	}

	var mutex sync.Mutex
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func() {
			for person := range peopleChan {
				pageRef, exists := personPageMap[person.FullName]

				if exists {
					texts := strings.Split(pageRef.Revision.Text, "\n")

					s := 0
					e := 0
					isPrevQuote := false
					for i, text := range texts {
						if len(text) > 2 && (text[0:2] == "*[" || text[0:3] == "* [") {
							isPrevQuote = false
							continue
						}

						if len(text) > 1 && text[0] == '*' && text[1] != '*' {
							s++
							isPrevQuote = true
						} else {
							if isPrevQuote && len(text) > 1 && text[0:2] == "**" {
								println(text)
								println(texts[i-1])
								e++
							}
							isPrevQuote = false
						}
					}

					mutex.Lock()
					cnt += s
					quoteSub += e
					mutex.Unlock()
				}
			}

			wg.Done()
		}()
	}

	close(peopleChan)
	wg.Wait()
	println(cnt)
	println(quoteSub)
	return
}