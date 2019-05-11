package repository

import (
	"github.com/gaaon/quote-collector/pkg/model"
	"github.com/gaaon/quote-collector/pkg/service/collect"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestWriter struct {
	content []byte
}

func (w *TestWriter) Write(p []byte) (n int, err error){
	w.content = append(w.content, p...)

	return 0, nil
}

func (w *TestWriter) ContentString() string {
	return string(w.content)
}

func TestSaveIntoWriter(t *testing.T) {
	assertT := assert.New(t)

	testWriter := &TestWriter{}
	var peopleList []model.Person

	peopleList = append(peopleList, model.Person{
		FullName: "Albert Einstein",
		ReversedName: "Einstein, Albert",
		Link: "/wiki/Albert_Einstein",
	})

	peopleList = append(peopleList, model.Person{
		FullName: "John von Neumann",
		ReversedName: "Neumann, John von",
		Link: "/wiki/John_von_Neumann",
	})

	err := collect.savePeopleListIntoWriter(testWriter, peopleList)
	assertT.NoError(err)

	assertT.Equal(
		testWriter.ContentString(),
		"Albert Einstein\tEinstein, Albert\t/wiki/Albert_Einstein\n" +
			"John von Neumann\tNeumann, John von\t/wiki/John_von_Neumann")
}

