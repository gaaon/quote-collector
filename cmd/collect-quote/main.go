package main

import (
	"github.com/gaaon/quote-collector/pkg/constant"
	"github.com/gaaon/quote-collector/pkg/model"
	"github.com/gaaon/quote-collector/pkg/repository"
	"github.com/gaaon/quote-collector/pkg/service/translate"
	"log"
	"os"
	"time"
)

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

		var (
			translatedByKakao string
			translatedByNaver string
			translatedByGoogle string
		)

		for _, quoteEntity := range quoteEntities {
			if len(quoteEntity.Content) > 100 {
				continue
			}

			if translatedByKakao, err = translate.TranslateByKakao(quoteEntity.Content); err != nil {
				log.Fatal(err)
			}

			if translatedByGoogle, err = translate.TranslateByGoogle(quoteEntity.Content); err != nil {
				log.Fatal(err)
			}

			if translatedByNaver, err = translate.TranslateByNaver(quoteEntity.Content); err != nil {
				log.Fatal(err)
			}

			if _, err = repository.InsertTranslation(model.TranslationEntity{
				Content: translatedByKakao,
				Vendor: constant.KAKAO,
				QuoteId: quoteEntity.Id,
			}); err != nil {
				log.Fatal(err)
			}

			if _, err = repository.InsertTranslation(model.TranslationEntity{
				Content: translatedByNaver,
				Vendor: constant.NAVER,
				QuoteId: quoteEntity.Id,
			}); err != nil {
				log.Fatal(err)
			}

			if _, err = repository.InsertTranslation(model.TranslationEntity{
				Content: translatedByGoogle,
				Vendor: constant.GOOGLE,
				QuoteId: quoteEntity.Id,
			}); err != nil {
				log.Fatal(err)
			}

			println("origin: ", quoteEntity.Content)
			println("translated(kakao): ", translatedByKakao)
			println("translated(naver): ", translatedByNaver)
			println("translated(google): ", translatedByGoogle)
			time.Sleep(10 * time.Second)
		}
	}
	}
}
