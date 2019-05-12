package main

import (
	"github.com/gaaon/quote-collector/pkg/model"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func (mainApp *MainApp) findKoreanNamesFromFullNames(peopleList []model.Person) error {

	f, _ := os.OpenFile("data/korean_snapshot.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	failed, _ := os.OpenFile("data/failed_to_find.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	for _, person := range peopleList {
		k, err := mainApp.nameTranslateService.TranslateFullNameToKorean(person.FullName)
		if err != nil {
			return err
		}

		if k == "" {
			_, _ = failed.WriteString(person.FullName + "\t" + person.ReversedName + "\n")
		} else {
			_, err = f.WriteString(person.FullName + "\t" + person.ReversedName + "\t" + k + "\n")
			if err != nil {
				return err
			}

			if err = mainApp.saveLastSuccessKoreanTranslation(person.FullName); err != nil {
				return err
			}
		}

		time.Sleep(time.Duration(mainApp.koreanTransIntervalInSec) * time.Second)
	}

	return nil
}

func (mainApp *MainApp) findLastSuccessKoreanTranslation(peopleList []model.Person) int {
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

func (mainApp *MainApp) saveLastSuccessKoreanTranslation(fullName string) error {
	return ioutil.WriteFile(
		"data/lastSuccessKoreanTrans.txt",
		[]byte(fullName),
		0644)
}