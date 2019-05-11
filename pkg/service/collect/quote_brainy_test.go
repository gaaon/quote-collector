package collect

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestFindVidAndPersonIdInBrainyByLinkPath(t *testing.T) {
	assertT := assert.New(t)

	path := "/authors/a_a_milne"
	vid, pid, err := FindVidAndPersonIdInBrainy(path)
	assertT.NoError(err)
	assertT.NotEmpty(vid)
	assertT.NotEmpty(pid)

	_, err = strconv.Atoi(pid)
	assertT.NoError(err)
}

func TestFindQuotesInBrainy(t *testing.T) {
	assertT := assert.New(t)

	path := "/authors/a_a_milne"
	vid, pid, err := FindVidAndPersonIdInBrainy(path)
	assertT.NoError(err)

	quotes, err := FindQuotesInBrainyWithPagination(vid, pid, 1)
	assertT.NoError(err)

	assertT.True(len(quotes) > 0)
}

func TestFindQuotesInBrainyByPath(t *testing.T) {
	assertT := assert.New(t)

	path := "/authors/a_a_milne"
	quotes, lastPagi, err := FindQuotesInBrainyByPath(path)
	assertT.NoError(err)

	assertT.True(lastPagi > 1)
	assertT.NotNil(quotes)

	path = "/authors/albert_einstein"
	quotes, lastPagi, err = FindQuotesInBrainyByPath(path)
	assertT.NoError(err)

	assertT.True(lastPagi > 5)
	assertT.NotNil(quotes)
	assertT.True(len(quotes) > 100)
}