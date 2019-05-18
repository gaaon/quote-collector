package main

import (
	"github.com/gaaon/quote-collector/pkg/repository"
	"github.com/gaaon/quote-collector/pkg/service/collect"
	"github.com/gaaon/quote-collector/pkg/service/notification"
	"github.com/gaaon/quote-collector/pkg/service/translate"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

type MainApp struct {
	task            string
	httpClient      *http.Client
	notiService 	notification.NotiService
	quoteBrainyService *collect.QuoteBrainyService
}

func NewMainApp(task string) (*MainApp, error) {
	httpClient := &http.Client{}

	notiService, err := notification.NewSlackService(httpClient)
	if err != nil {
		return nil, err
	}

	return &MainApp{
		task:       task,
		httpClient: httpClient,
		notiService: notiService,
		quoteBrainyService: collect.NewQuoteBrainyService(httpClient),
	}, nil
}

func (mainApp *MainApp) Run() {
	switch mainApp.task {
	case "find":
		{
			var err error
			if err = mainApp.notiService.SendNotiToDevice("[quotes] start collecting quotes from brainy"); err != nil {
				log.Fatal(err)
			}

			if err = mainApp.findQuotesFromBrainy(); err != nil {
				log.Error(err)
				err2 := mainApp.notiService.SendNotiToDevice("[quotes] collecting quotes has problem\n" + err.Error())
				if err2 != nil {
					log.Fatal(err2)
				}
			}
		}
	case "translate":
		{
			quoteEntities, err := repository.FindQuoteEntitiesFromDB()
			if err != nil {
				log.Fatal(err)
			}

			startIdx := mainApp.findLastSuccessQuoteTranslation(quoteEntities) + 1
			println("find startIdx: ", startIdx)

			for i, quoteEntity := range quoteEntities {
				if i < startIdx {
					continue
				}

				content := mainApp.filterQuoteContent(quoteEntity.Content)
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

				if err = mainApp.saveLastSuccessQuoteTranslation(quoteEntity.Id); err != nil {
					log.Fatal(err)
				}

				time.Sleep(60 * time.Second)
			}
		}
	}
}

func main() {
	log.Info("Start collecting quotes")
	task, exists := os.LookupEnv("COLLECT_TASK")
	if !exists {
		task = "find"
	}

	mainApp, err := NewMainApp(task)
	if err != nil {
		log.Fatal(err)
	}

	mainApp.Run()
}
