package quotewiki

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestGetPeopleListHtmlByA(t *testing.T) {
	assertT := assert.New(t)

	bodyReader, err := getPeopleListHtmlByName("A")
	assertT.NoError(err)
	assertT.NotNil(bodyReader)
}

func TestGetPeopleListByA(t *testing.T) {
	assertT := assert.New(t)

	bodyReader, err := getPeopleListHtmlByName("A")
	assertT.NoError(err)

	people, err := GetPeopleListByReader(bodyReader)
	assertT.NoError(err)
	assertT.NotNil(people)
	assertT.True(len(people) > 0)
}

func TestGetPeopleListFromAToZ(t *testing.T) {
	assertT := assert.New(t)

	people, err := GetPeopleListFromAToZ()
	assertT.NoError(err)
	assertT.NotNil(people)

	firstPeople := people[0]
	assertT.True(strings.ToLower(string(firstPeople[0])) == "a")

	lastPeople := people[len(people) - 1]
	assertT.True(strings.ToLower(string(lastPeople[0])) == "z")
}