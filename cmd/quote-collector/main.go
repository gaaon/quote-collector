package main

import (
	"fmt"
	"github.com/gaaon/quote-collector/pkg/google"
	"github.com/gaaon/quote-collector/pkg/quotewiki"
	"log"
	"os"
	"time"
)

type CompositeName struct {
	Original string
	Korean string
}

type CompositeNames []CompositeName

func (c CompositeNames) Len() int {
	return len(c)
}

func (c CompositeNames) Less(i, j int) bool {
	return c[i].Original <c[j].Original
}

func (c CompositeNames) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func main() {
	peopleList, err := quotewiki.GetPeopleListFromSnapshot()
	if err != nil {
		log.Fatal(err)
	}

	f, _ := os.Create("composite_snapshot.txt")
	failed, _ := os.Create("failed_to_find.txt")
	for i := 0; i < len(peopleList); i++ {
		original := peopleList[i]
		k, err  := google.GetKoreanNameFromEnglish(original)
		if err != nil {
			fmt.Println(err.Error())
		}

		if k == "" {
			_, _ = failed.WriteString(original + "\n")
		}
		_, _ = f.WriteString(original + "\t" + k + "\n")

		if i % 100 == 0 {
			fmt.Printf("%d개 다운 성공\n", i)
		}
		time.Sleep(10 * time.Second)
	}
}
