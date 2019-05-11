package collect

import (
	"github.com/stretchr/testify/assert"
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
}