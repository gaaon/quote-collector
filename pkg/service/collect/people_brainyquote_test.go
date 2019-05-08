package collect

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindPeopleListStartsWithA(t *testing.T) {
	assertT := assert.New(t)

	peopleList, err := FindPeopleListInBrainyStartsWith("a", 2)
	assertT.NoError(err)
	assertT.True(len(peopleList) > 0)
}

func TestFindPeopleListWithPagination(t *testing.T) {
	assertT := assert.New(t)

	_, err := FindPeopleListInBrainyStartsWith("a", 13)
	assertT.Errorf(err, "page not exists")

	_, err = FindPeopleListInBrainyStartsWith("a", 100)
	assertT.Errorf(err, "page not exists")
}

