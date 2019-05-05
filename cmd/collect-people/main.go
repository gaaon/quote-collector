package main

import (
	"fmt"
	"github.com/gaaon/quote-collector/pkg/google"
	"github.com/gaaon/quote-collector/pkg/quotewiki"
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

func findKoreanNameFromEng(peopleList []quotewiki.Person) {
	f, _ := os.Create("data/korean_snapshot.txt")
	failed, _ := os.Create("data/failed_to_find.txt")

	for i := 0; i < len(peopleList); i++ {
		original := peopleList[i]
		k, err  := google.GetKoreanNameFromEnglish(original.FullName)
		if err != nil {
			fmt.Println(err.Error())
		}

		if k == "" {
			newName := original.ReversedName
			k, err = google.GetKoreanNameFromEnglish(newName)
			println("[newName, newK] ", newName, k)
			if err != nil {
				fmt.Println(err.Error())
			}
		}

		if k == "" {
			_, _ = failed.WriteString(original.FullName + "\t" + original.ReversedName + "\n")
		}
		_, _ = f.WriteString(original.FullName + "\t" + original.ReversedName + "\t" + k + "\n")

		if i % 100 == 0 {
			fmt.Printf("%d개 다운 성공\n", i)
		}

		time.Sleep(10 * time.Second)
	}
}

func migrateKoreanSnapshotWithDefault(peopleList []quotewiki.Person, koreanNameMap map[string]string) (
	migratedPersonList []quotewiki.Person, err error) {

	for _, person := range peopleList {
		koreanName, exists := koreanNameMap[person.FullName]
		if exists {
			migratedPersonList = append(migratedPersonList, quotewiki.Person{
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

	switch env {
	case "file": {
		switch task {
		case "find": {
			peopleList, err := quotewiki.FindPeopleListFromSnapshot()
			if err != nil {
				log.Fatal(err)
			}

			println("Find people count ", len(peopleList))
		}
		case "korean": {
			peopleList, err := quotewiki.FindPeopleListFromSnapshot()
			if err != nil {
				log.Fatal(err)
			}

			findKoreanNameFromEng(peopleList)
		}
		default:
			log.Fatal("no such collect task")
		}

	}

	case "db": {
		switch task {
		case "migrate": {
			peopleList, err := quotewiki.FindPeopleListFromSnapshot()
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
				_, err := quotewiki.InsertPersonIntoDB(person.FullName, person.KoreanName, person.Link)
				if err != nil {
					log.Fatal(nil)
				}
			}
		}
		case "find": {
			peopleList, err := quotewiki.FindPeopleListFromDB()
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
