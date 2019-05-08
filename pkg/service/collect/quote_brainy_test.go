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

	println(vid, pid)
}