package repository

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type DBTestSuite struct {
	suite.Suite
}

func (suite *DBTestSuite) SetupSuite() {
	SetDBName("test-quote-collector")
}

func (suite *DBTestSuite) SetupTest() {
	err := dropPeopleCollection()
	suite.NoError(err)
}

func (suite *DBTestSuite) TearDownSuite() {
	RecoverDBName()
}

func TestDBTestSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}

func (suite *DBTestSuite) TestFindPeopleList() {
	assertT := assert.New(suite.T())

	_, err := InsertPerson("Zhuge Liang", "제갈량", "")
	assertT.NoError(err)

	peopleList, err := FindPeopleList()
	assertT.NoError(err)

	assertT.Len(peopleList, 1)
}
