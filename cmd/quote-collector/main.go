package main

import (
	"fmt"
	"github.com/gaaon/quote-collector/pkg/google"
	"github.com/gaaon/quote-collector/pkg/quotewiki"
	"log"
	"os"
	"time"
)

func getFileVersion() (fileVersion string) {
	if len(os.Args) > 1 {
		fileVersion = os.Args[1]
	} else {
		fileVersion = "1"
	}

	return
}

func main() {
	fileVersion := getFileVersion()
	_, err := os.Stat("data/" + fileVersion)
	if os.IsNotExist(err) {
		_ = os.MkdirAll("data/" + fileVersion, os.ModePerm)
	}

	peopleList, err := quotewiki.GetPeopleListFromSnapshot(fileVersion)
	if err != nil {
		log.Fatal(err)
	}

	f, _ := os.Create("data/" + fileVersion + "/composite_snapshot.txt")
	failed, _ := os.Create("data/" + fileVersion + "/failed_to_find.txt")
	for i := 0; i < len(peopleList); i++ {
		original := peopleList[i]
		k, err  := google.GetKoreanNameFromEnglish(original.TitleName)
		if err != nil {
			fmt.Println(err.Error())
		}

		if k == "" {
			newName := original.TextName
			k, err = google.GetKoreanNameFromEnglish(newName)
			println("[newName, newK] ", newName, k)
			if err != nil {
				fmt.Println(err.Error())
			}
		}

		if k == "" {
			_, _ = failed.WriteString(original.TitleName + "\t" + original.TextName + "\n")
		}
		_, _ = f.WriteString(original.TitleName + "\t" + original.TextName + "\t" + k + "\n")

		if i % 100 == 0 {
			fmt.Printf("%d개 다운 성공\n", i)
		}
		time.Sleep(10 * time.Second)
	}
}
