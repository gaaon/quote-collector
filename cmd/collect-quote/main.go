package main

import (
	"github.com/gaaon/quote-collector/pkg/model"
	"github.com/gaaon/quote-collector/pkg/repository"
	"github.com/gaaon/quote-collector/pkg/service/translate"
	"log"
	"os"
	"strings"
	"time"
)

func filterQuoteContent(content string) string{
	content = strings.ReplaceAll(content, "<br/>", " ")
	content = strings.ReplaceAll(content, "<br>", " ")
	content = strings.ReplaceAll(content, "<BR>", " ")
	return strings.ReplaceAll(content, "<br />", " ")
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

func main() {
	task, exists := os.LookupEnv("COLLECT_TASK")
	if !exists {
		task = "find"
	}

	switch task {
	case "find": {
		findQuotesFromMediaWiki()
	}
	case "translate": {
		quoteEntities, err := repository.FindQuoteEntitiesFromDB()
		if err != nil {
			log.Fatal(err)
		}

		for _, quoteEntity := range quoteEntities {
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
			time.Sleep(10 * time.Second)
		}
	}
	}
}
