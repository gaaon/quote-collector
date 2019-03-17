package main

import (
	"fmt"
	"github.com/gaaon/quote-collector/pkg/google"
	"github.com/gaaon/quote-collector/pkg/quotewiki"
	"log"
	"sort"
	"sync"
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

	nameChan := make(chan string, len(peopleList))

	for _, name := range peopleList[0: 1] {
		nameChan <- name
	}

	var names []CompositeName
	var mutex sync.Mutex
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(num int) {
			for name := range nameChan {
				k, err  := google.GetKoreanNameFromEnglish(name)
				if err != nil {
					fmt.Println(err.Error())
				}
				mutex.Lock()
				names = append(names, CompositeName{name, k})
				fmt.Println(CompositeName{name, k})
				mutex.Unlock()
			}

			wg.Done()
		}(i)
	}

	close(nameChan)
	wg.Wait()

	sort.Sort(CompositeNames(names))
	for _, cName := range names {
		println(cName.Original, cName.Korean)
	}
}
