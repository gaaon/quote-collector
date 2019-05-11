package main

import (
	"fmt"
	"github.com/gaaon/quote-collector/pkg/model"
	"github.com/gaaon/quote-collector/pkg/repository"
	"github.com/gaaon/quote-collector/pkg/service/collect"
	"github.com/gaaon/quote-collector/pkg/service/notification"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

func findKoreanNameMapFromSnapshot() (koreanNameMap map[string]string, err error){
	koreanNameMap = make(map[string]string)

	var f *os.File
	f, err = os.Open("data/korean_snapshot.txt")
	if err != nil {
		return
	}

	defer f.Close()

	var content []byte
	if content, err = ioutil.ReadAll(f); err != nil {
		return
	}

	splits := strings.Split(string(content), "\n")
	for _, split := range splits {
		if split == "" {
			continue
		}

		info := strings.Split(split, "\t")

		koreanNameMap[info[0]] = info[2]
	}

	return
}

var sendNoti = false

func findKoreanNameFromEng(nameTranslateService *collect.NameTranslateService, peopleList []model.Person) {
	f, _ := os.OpenFile("data/korean_snapshot.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	failed, _ := os.OpenFile("data/failed_to_find.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	for i := 0; i < len(peopleList); i++ {
		original := peopleList[i]
		k, err  := nameTranslateService.TranslateFullNameToKorean(original.FullName)
		if err != nil {
			fmt.Println(err.Error())

			if err.Error() == "too many request status code from server" && sendNoti == false {
				if err2 := notification.SendNotiToDevice("429 comes", "quote-collector server"); err2 != nil {
					log.Fatal(err2)
				}
				sendNoti = true
			}

			continue
		}

		if k == "" {
			_, _ = failed.WriteString(original.FullName + "\t" + original.ReversedName + "\n")
		} else {
			_, err = f.WriteString(original.FullName + "\t" + original.ReversedName + "\t" + k + "\n")
			if err != nil {
				log.Fatal(err)
			}
			if err = saveLastSuccessKoreanTranslation(original.FullName); err != nil {
				log.Fatal(err)
			}
		}

		if i % 100 == 0 {
			fmt.Printf("%d개 다운 성공\n", i)
		}

		interval := 60
		println("sleep time: ", interval, "seconds")
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func findLastSuccessKoreanTranslation(peopleList []model.Person) int {
	f, err := os.Open("data/lastSuccessKoreanTrans.txt")
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

func saveLastSuccessKoreanTranslation(fullName string) error {
	return ioutil.WriteFile(
		"data/lastSuccessKoreanTrans.txt",
		[]byte(fullName),
		0644)
}

func migrateKoreanSnapshotWithDefault(peopleList []model.Person, koreanNameMap map[string]string) (
	migratedPersonList []model.Person, err error) {

	for _, person := range peopleList {
		koreanName, exists := koreanNameMap[person.FullName]
		if exists {
			migratedPersonList = append(migratedPersonList, model.Person{
				FullName: person.FullName,
				ReversedName: person.ReversedName,
				Link: person.Link,
				KoreanName: koreanName,
			})
		}
	}

	return
}

func main() {
	println("Start collecting people list")

	env, exists := os.LookupEnv("COLLECT_ENV")
	if !exists {
		env = "file"
	}

	task, exists := os.LookupEnv("COLLECT_TASK")
	if !exists {
		task = "find"
	}

	peopleSnapshotService := collect.NewPeopleSnapshotService()
	nameTranslateService := collect.NewNameTranslateService()
	brainyQuoteService, err := collect.NewBrainyQuoteService(peopleSnapshotService)
	if err != nil {
		log.Fatal(err)
	}

	switch env {
	case "file": {
		switch task {
		case "find": {
			peopleList, err := brainyQuoteService.FindPeopleListFromSnapshot()
			if err != nil {
				log.Fatal(err)
			}

			println("Find people count ", len(peopleList))
		}
		case "korean": {
			peopleList, err := brainyQuoteService.FindPeopleListFromSnapshot()
			if err != nil {
				log.Fatal(err)
			}

			hoursToCollect := len(peopleList) * 60 / 60 / 60
			lastIndex := findLastSuccessKoreanTranslation(peopleList)

			fmt.Printf("time for finding: %d hours\n", hoursToCollect)

			if err = notification.SendNotiToDevice("test start message", "quote-collector server"); err != nil {
				log.Fatal(err)
			}
			findKoreanNameFromEng(nameTranslateService, peopleList[lastIndex+1:])
		}
		default:
			log.Fatal("no such collect task")
		}

	}

	case "db": {
		switch task {
		case "migrate": {
			peopleList, err := peopleSnapshotService.FindPeopleList()
			if err != nil {
				log.Fatal(err)
			}

			if err = repository.InsertPeopleListIntoDB(peopleList); err != nil {
				log.Fatal(err)
			}
		}
		case "find": {
			peopleList, err := repository.FindPeopleList()
			if err != nil {
				log.Fatal(err)
			}

			println("Find people count ", len(peopleList))
		}
		}
	}
	default:
		log.Fatal("no such collect env")
	}

}
