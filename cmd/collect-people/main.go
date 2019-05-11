package main

import (
	"fmt"
	"github.com/gaaon/quote-collector/pkg/model"
	"github.com/gaaon/quote-collector/pkg/repository"
	"github.com/gaaon/quote-collector/pkg/service/collect"
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
			fmt.Printf("time for finding: %d hours\n", hoursToCollect)

			f, _ := os.Create("data/korean_snapshot.txt")
			failed, _ := os.Create("data/failed_to_find.txt")

			for i := len(peopleList) - 1; i >= 0; i-- {
				original := peopleList[i]
				k, err  := nameTranslateService.TranslateFullNameToKorean(original.FullName)
				if err != nil {
					fmt.Println(err.Error())
				}

				if k == "" {
					_, _ = failed.WriteString(original.FullName + "\t" + original.ReversedName + "\n")
				}
				_, _ = f.WriteString(original.FullName + "\t" + original.ReversedName + "\t" + k + "\n")

				if i % 100 == 0 {
					fmt.Printf("%d개 다운 성공\n", i)
				}

				time.Sleep(60 * time.Second)
			}
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

			koreanNameMap, err := findKoreanNameMapFromSnapshot()
			if err != nil {
				log.Fatal(err)
			}

			migratedPeopleList, err := migrateKoreanSnapshotWithDefault(peopleList, koreanNameMap)
			if err != nil {
				log.Fatal(err)
			}

			for _, person := range migratedPeopleList {
				_, err := repository.InsertPerson(person.FullName, person.KoreanName, person.Link)
				if err != nil {
					log.Fatal(nil)
				}
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
