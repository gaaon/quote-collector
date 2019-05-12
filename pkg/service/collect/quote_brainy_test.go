package collect

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"strconv"
	"testing"
)

type QuoteBrainyServiceTestSuite struct {
	suite.Suite
	service *quoteBrainyService
}

func (suite *QuoteBrainyServiceTestSuite) SetupTest() {
	suite.service = NewQuoteBrainyService()
}

func TestQuoteBrainyServiceTestSuite(t *testing.T) {
	suite.Run(t, new(QuoteBrainyServiceTestSuite))
}


func (suite *QuoteBrainyServiceTestSuite) TestFindVidAndPersonIdByLink() {
	assertT := assert.New(suite.T())

	path := "/authors/a_a_milne"
	vid, pid, err := suite.service.findVidAndPersonId(path)
	assertT.NoError(err)
	assertT.NotEmpty(vid)
	assertT.NotEmpty(pid)

	_, err = strconv.Atoi(pid)
	assertT.NoError(err)
}

func (suite *QuoteBrainyServiceTestSuite) TestFindQuotesInBrainyByLink() {
	assertT := assert.New(suite.T())

	link := "/authors/a_a_milne"
	vid, pid, err := suite.service.findVidAndPersonId(link)
	assertT.NoError(err)

	quotes, err := suite.service.findQuotesWithPagination(vid, pid, 1)
	assertT.NoError(err)

	assertT.True(len(quotes) > 0)
}

func (suite *QuoteBrainyServiceTestSuite) TestFindSequentialQuotesByLink() {
	assertT := assert.New(suite.T())

	link := "/authors/a_a_milne"
	quotes, lastPagi, err := suite.service.FindAllQuotesByLink(link)
	assertT.NoError(err)

	assertT.True(lastPagi > 1)
	assertT.NotNil(quotes)

	link = "/authors/albert_einstein"
	quotes, lastPagi, err = suite.service.FindAllQuotesByLink(link)
	assertT.NoError(err)

	assertT.True(lastPagi > 5)
	assertT.NotNil(quotes)
	assertT.True(len(quotes) > 100)
}