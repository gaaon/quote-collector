package main

import (
	"github.com/gaaon/quote-collector/pkg/model"
	"github.com/gaaon/quote-collector/pkg/repository"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func (mainApp *MainApp) filterQuoteContent(content string) string {
	replacer := strings.NewReplacer(
		"<br/>", " ",
		"<br>", " ",
		"<BR>", " ",
		"<br />", " ",
		"&quot;", "'",
		"&lt;", " ",
		"&gt;", " ")

	return replacer.Replace(content)
}

func (mainApp *MainApp) findQuotesFromBrainy() error {
	peopleList, err := repository.FindPeopleList()
	if err != nil {
		return err
	}

	lastSuccessIdx := mainApp.findLastSuccessQuoteCollect(peopleList)

	for i, person := range peopleList {
		if i == lastSuccessIdx {
			log.Info("continue collecting quotes after %s", person.FullName)
		}
		if i <= lastSuccessIdx {
			continue
		}

		partialQuotes, _, err := mainApp.quoteBrainyService.FindAllQuotesByLink(person.Link)
		if err != nil {
			return err
		}

		quoteEntities := repository.GetQuoteEntitiesWithPerson(partialQuotes, person)
		if err = repository.InsertQuoteEntitiesIntoDB(quoteEntities); err != nil {
			return err
		}

		if err = mainApp.saveLastSuccessQuoteCollect(person.FullName); err != nil {
			return err
		}

		time.Sleep(10 * time.Second)
	}

	return nil
}

func (mainApp *MainApp) findLastSuccessQuoteCollect(peopleList []model.Person) int {
	f, err := os.Open("data/lastSuccessQuoteCollect.txt")
	if os.IsNotExist(err) {
		return -1
	}
	defer f.Close()

	contentRaw, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	content := string(contentRaw)

	if content == "" {
		return -1
	} else {
		for i, person := range peopleList {
			if person.FullName == content {
				return i
			}
		}

		return -1
	}
}

func (mainApp *MainApp) saveLastSuccessQuoteCollect(fullName string) error {
	return ioutil.WriteFile(
		"data/lastSuccessQuoteCollect.txt",
		[]byte(fullName),
		0644)
}

func (mainApp *MainApp) findLastSuccessQuoteTranslation(entities []model.QuoteEntity) int {

	f, err := os.Open("data/lastSuccessQuoteTrans.txt")
	if os.IsNotExist(err) {
		return -1
	}
	defer f.Close()

	contentRaw, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	content := string(contentRaw)

	if content == "" {
		return -1
	} else {
		for i, quote := range entities {
			if quote.Id.String() == content {
				return i
			}
		}

		return -1
	}
}

func (mainApp *MainApp) saveLastSuccessQuoteTranslation(id *primitive.ObjectID) error {
	return ioutil.WriteFile(
		"data/lastSuccessQuoteTrans.txt",
		[]byte(id.String()),
		0644)
}
