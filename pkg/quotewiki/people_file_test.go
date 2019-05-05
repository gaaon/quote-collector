package quotewiki

import (
	"github.com/stretchr/testify/assert"
	"strings"
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
	var peopleList []Person

	peopleList = append(peopleList, Person{
		FullName: "Albert Einstein",
		ReversedName: "Einstein, Albert",
		Link: "/wiki/Albert_Einstein",
	})

	peopleList = append(peopleList, Person{
		FullName: "John von Neumann",
		ReversedName: "Neumann, John von",
		Link: "/wiki/John_von_Neumann",
	})

	err := saveIntoWriter(peopleList, testWriter)
	assertT.NoError(err)

	assertT.Equal(
		testWriter.ContentString(),
		"Albert Einstein\tEinstein, Albert\t/wiki/Albert_Einstein\n" +
			"John von Neumann\tNeumann, John von\t/wiki/John_von_Neumann")
}

func TestGetPeopleListHtmlByA(t *testing.T) {
	assertT := assert.New(t)

	bodyReader, err := getPeopleListHtmlByName("A")
	assertT.NoError(err)

	defer bodyReader.Close()
	assertT.NotNil(bodyReader)
}

func TestGetPeopleListFromAToZ(t *testing.T) {
	assertT := assert.New(t)

	people, err := getPeopleListFromAToZ()
	assertT.NoError(err)
	assertT.NotNil(people)

	firstPeople := people[0]
	assertT.True(strings.ToLower(string(firstPeople.FullName[0])) == "a")

	lastPeople := people[len(people) - 1]
	assertT.True(strings.ToLower(string(lastPeople.FullName[0])) == "z")
}