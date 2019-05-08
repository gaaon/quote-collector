package main

import (
	"github.com/gaaon/quote-collector/pkg/model"
	"github.com/gaaon/quote-collector/pkg/repository"
	"github.com/gaaon/quote-collector/pkg/service/translate"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

func filterQuoteContent(content string) string{
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

func findQuotesFromMediaWiki() {
	peopleList, err := repository.FindPeopleList()
	if err != nil {
		log.Fatal(err)
	}

	mediaWikiXmlFile, err := os.Open("data/enwikiquote-latest-pages-articles.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer mediaWikiXmlFile.Close()

	mediaWiki, err := repository.GetMediaWikiFromReader(mediaWikiXmlFile)
	if err != nil {
		log.Fatal(err)
	}

	pageMap := repository.GetPersonNamePageMapFromMediaWiki(mediaWiki)

	var quoteEntities []model.QuoteEntity
	for _, person := range peopleList {
		partialQuotes, err := repository.FindQuotesInPageMapByFullName(pageMap, person.FullName)
		if err != nil {
			log.Println(err)
		}

		quoteEntities = append(quoteEntities, repository.GetQuoteEntitiesWithPerson(partialQuotes, person)...)
	}

	if err = repository.InsertQuoteEntitiesIntoDB(quoteEntities); err != nil {
		log.Fatal(err)
	}

	total := 0
	for _, quoteEntity := range quoteEntities {
		total += len(quoteEntity.Content)
	}

	println("total quotes count: ", len(quoteEntities))
	println("total characters count: ", total)
}

func findLastSuccessQuoteTranslation(entities []model.QuoteEntity) int {
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

func saveLastSuccessQuoteTranslation(id *primitive.ObjectID) error {
	return ioutil.WriteFile(
		"data/lastSuccessQuoteTrans.txt",
		[]byte(id.String()),
		0644)
}

func main() {
	task, exists := os.LookupEnv("COLLECT_TASK")
	if !exists {
		task = "find"
	}

	switch task {
	case "find": {
		//findQuotesFromMediaWiki()
	}
	case "translate": {
		quoteEntities, err := repository.FindQuoteEntitiesFromDB()
		if err != nil {
			log.Fatal(err)
		}

		startIdx := findLastSuccessQuoteTranslation(quoteEntities) + 1
		println("find startIdx: ", startIdx)

		for i, quoteEntity := range quoteEntities {
			if i < startIdx {
				continue
			}

			content := filterQuoteContent(quoteEntity.Content)
			if len(content) > 100 {
				continue
			}

			translatedByNaver, _, err := translate.FindTranslationByNaverAndSave(content, quoteEntity)
			if err != nil {
				log.Fatal(err)
			}

			translatedByGoogle, _, err := translate.FindTranslationByGoogleAndSave(content, quoteEntity)
			if err != nil {
				log.Fatal(err)
			}

			translatedByKakao, _, err := translate.FindTranslationByKakaoAndSave(content, quoteEntity)
			if err != nil {
				log.Fatal(err)
			}

			println("origin: ", content)
			println("translated(kakao): ", translatedByKakao)
			println("translated(naver): ", translatedByNaver)
			println("translated(google): ", translatedByGoogle)

			if err = saveLastSuccessQuoteTranslation(quoteEntity.Id); err != nil {
				log.Fatal(err)
			}

			time.Sleep(10 * time.Second)
		}
	}
	}
}
