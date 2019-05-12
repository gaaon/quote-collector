package collect

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type BrainyQuoteServiceTestSuite struct {
	suite.Suite
	service *peopleBrainyService
}

func (suite *BrainyQuoteServiceTestSuite) SetupTest() {
	var (
		err error
		peopleSnapshotService = NewPeopleSnapshotService()
	)

	suite.service, err = NewBrainyQuoteService(peopleSnapshotService)
	assert.NoError(suite.T(), err)
}

func TestBrainyQuoteServiceTestSuite(t *testing.T) {
	suite.Run(t, new(BrainyQuoteServiceTestSuite))
}

func (suite *BrainyQuoteServiceTestSuite) TestFindPeopleListStartsWithA() {
	assertT := assert.New(suite.T())

	peopleList, err := suite.service.findPeopleListStartsWith("a", 2)
	assertT.NoError(err)
	assertT.True(len(peopleList) > 0)
}

func (suite *BrainyQuoteServiceTestSuite) TestFindPeopleListWithPagination() {
	assertT := assert.New(suite.T())

	peopleList, err := suite.service.findPeopleListStartsWith("a", 13)
	assertT.NoError(err)
	assertT.Len(peopleList, 0)

	peopleList, err = suite.service.findPeopleListStartsWith("a", 100)
	assertT.NoError(err)
	assertT.Len(peopleList, 0)
}

