package collect

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestGetPeopleListHtmlByA(t *testing.T) {
	assertT := assert.New(t)

	bodyReader, err := FindPeopleListHtmlByName("A")
	assertT.NoError(err)
	defer bodyReader.Close()

	assertT.NotNil(bodyReader)
}

func TestGetPeopleListFromAToZ(t *testing.T) {
	assertT := assert.New(t)

	people, err := FindPeopleListFromAToZ()
	assertT.NoError(err)
	assertT.NotNil(people)

	firstPeople := people[0]
	assertT.True(strings.ToLower(string(firstPeople.FullName[0])) == "a")

	lastPeople := people[len(people) - 1]
	assertT.True(strings.ToLower(string(lastPeople.FullName[0])) == "z")
}