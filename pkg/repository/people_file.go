package repository

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gaaon/quote-collector/pkg/model"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

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
					FullName: name,
					ReversedName: nameLink.Text(),
					Link: linkUrl,
				})
			}
		})
	})

	return
}

func savePeopleListIntoWriter(writer io.Writer, peopleList []model.Person) (err error) {
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

const peopleListSnapshotLocation = "data/snapshot.txt"
func SavePeopleListIntoSnapshot(peopleList []model.Person) (err error) {
	f, err := os.Create(peopleListSnapshotLocation)
	if err != nil {
		return
	}
	defer f.Close()

	return savePeopleListIntoWriter(f, peopleList)
}

func IsExistPeopleListSnapshot() bool {
	f, err := os.Open(peopleListSnapshotLocation)
	defer f.Close()

	return os.IsExist(err)
}

func FindPeopleListFromSnapshot() (peopleList []model.Person, err error) {
	f, err := os.Open(peopleListSnapshotLocation)
	if err != nil {
		return
	}

	var content []byte
	if content, err = ioutil.ReadAll(f); err != nil {
		return
	}

	splits := strings.Split(string(content), "\n")
	for _, split := range splits {
		names := strings.Split(split, "\t")
		peopleList = append(peopleList, model.Person{FullName: names[0], ReversedName: names[1], Link: names[2]})
	}

	return
}