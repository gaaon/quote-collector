package quotewiki

import (
	"github.com/stretchr/testify/assert"
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

	people, err := GetPeopleListByName(bodyReader, "A")
	assertT.NoError(err)
	assertT.NotNil(people)
	assertT.True(len(people) > 0)
}