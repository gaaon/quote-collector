package main

import (
	"fmt"
	"github.com/gaaon/quote-collector/pkg/google"
	"github.com/gaaon/quote-collector/pkg/quotewiki"
	"log"
	"os"
	"time"
)

func getKoreanNameFromEng() {
	peopleList, err := quotewiki.FindPeopleListFromSnapshot()
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

func getQuotesFromPeopleList() {
	quotes, err := quotewiki.FindPeopleListFromSnapshot()
	if err != nil {
		log.Fatal(err)
	}

	for _, name := range quotes {
		println(name.TitleName)
	}
	//quotewiki.GetQuotesFromCompositeName(quotes[0])
}

func main() {
	peopleList, err := quotewiki.FindPeopleListFromDB()
	if err != nil {
		log.Fatal(err)
	}

	println(len(peopleList))
}
