package main

import (
	"github.com/gaaon/quote-collector/pkg/repository"
	"github.com/gaaon/quote-collector/pkg/service/collect"
	"github.com/gaaon/quote-collector/pkg/service/notification"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

type MainApp struct {
	httpClient               *http.Client
	env                      string
	task                     string
	koreanTransIntervalInSec int
	nameTranslateService     *collect.NameTranslateService
	peopleBrainyService      *collect.PeopleBrainyService
	notiService          	 notification.NotiService
}

func NewMainApp(env string, task string) (*MainApp, error) {
	client := &http.Client{}

	peopleSnapshotService := collect.NewPeopleSnapshotService()
	peopleBrainyService, err := collect.NewPeopleBrainyService(peopleSnapshotService)
	if err != nil {
		return nil, err
	}

	notiService, err := notification.NewSlackService(client)
	if err != nil {
		return nil, err
	}

	return &MainApp{
		httpClient:               client,
		env:                      env,
		task:                     task,
		koreanTransIntervalInSec: 60,
		nameTranslateService:     collect.NewNameTranslateService(client),
		peopleBrainyService:      peopleBrainyService,
		notiService:              notiService,
	}, nil
}

func (mainApp *MainApp) Run() {
	switch mainApp.env {
	case "file":
		{
			switch mainApp.task {
			case "find":
				{
					peopleList, err := mainApp.peopleBrainyService.FindPeopleListFromSnapshotOrRemote()
					if err != nil {
						log.Fatal(err)
					}

					log.Infof("find people count %d", len(peopleList))
				}
			case "korean":
				{
					log.Info("start translating korean names")

					peopleList, err := mainApp.peopleBrainyService.FindPeopleListFromSnapshotOrRemote()
					if err != nil {
						log.Fatal(err)
					}

					hoursToCollect := len(peopleList) * mainApp.koreanTransIntervalInSec / 60 / 60
					log.Infof("time for finding: %d hours\n", hoursToCollect)

					lastIndex := mainApp.findLastSuccessKoreanTranslation(peopleList)
					if err = mainApp.notiService.SendNotiToDevice("[korean] start translating fullname to korean"); err != nil {
						log.Fatal(err)
					}

					if err = mainApp.findKoreanNamesFromFullNames(peopleList[lastIndex+1:]); err != nil {
						log.Error(err)

						var err2 error
						if err.Error() == "too many request status code from server" {
							err2 = mainApp.notiService.SendNotiToDevice("[korean] 429 comes")
						} else {
							err2 = mainApp.notiService.SendNotiToDevice("[korean] something happens\n" + err.Error())
						}

						if err2 != nil {
							log.Fatal(err2)
						}
					} else {
						err2 := mainApp.notiService.SendNotiToDevice("[korean] success translating korean names")
						if err2 != nil {
							log.Fatal(err2)
						}
					}
				}
			default:
				log.Fatal("no such collect task")
			}
		}

	case "db":
		{
			switch mainApp.task {
			case "migrate":
				{
					peopleList, err := mainApp.peopleBrainyService.FindPeopleListFromSnapshotOrRemote()
					if err != nil {
						log.Fatal(err)
					}

					if err = repository.InsertPeopleListIntoDB(peopleList); err != nil {
						log.Fatal(err)
					}
				}
			case "find":
				{
					peopleList, err := repository.FindPeopleList()
					if err != nil {
						log.Fatal(err)
					}

					println("Find people count ", len(peopleList))
				}
			default:
				log.Fatal("no such collect task")
			}
		}
	default:
		log.Fatal("no such collect env")
	}
}

func main() {
	log.Info("Start collecting people list")

	env, exists := os.LookupEnv("COLLECT_ENV")
	if !exists {
		env = "file"
	}

	task, exists := os.LookupEnv("COLLECT_TASK")
	if !exists {
		task = "find"
	}

	mainApp, err := NewMainApp(env, task)
	if err != nil {
		log.Fatal(err)
	}
	mainApp.Run()
}
