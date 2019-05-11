package collect

import (
	"bytes"
	"github.com/gaaon/quote-collector/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type PeopleSnapshotServiceTestSuite struct {
	suite.Suite
	service *peopleSnapshotService
}

func (suite *PeopleSnapshotServiceTestSuite) SetupTest() {
	suite.service = NewPeopleSnapshotService()
}

func TestPeopleSnapshotServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PeopleSnapshotServiceTestSuite))
}

func (suite *PeopleSnapshotServiceTestSuite) TestWritePeopleListIntoWriter() {
	assertT := assert.New(suite.T())

	peopleList := []model.Person{
		{
			FullName: "Albert Einstein",
			Link: "/authors/albert_einstein",
		},
		{
			FullName: "Michael Jackson",
			Link: "/authors/michael_jackson",
		},
	}

	var buffer bytes.Buffer
	err := suite.service.writePeopleListIntoWriter(&buffer, peopleList)
	assertT.NoError(err)
	assertT.Equal(
		"Albert Einstein\t\t/authors/albert_einstein\nMichael Jackson\t\t/authors/michael_jackson",
		buffer.String())
}

func (suite *PeopleSnapshotServiceTestSuite) TestReadPeopleListFromReader() {
	assertT := assert.New(suite.T())

	var buffer bytes.Buffer
	_, err := buffer.WriteString(
		"Albert Einstein\t\t/authors/albert_einstein\nMichael Jackson\t\t/authors/michael_jackson")
	assertT.NoError(err)

	peopleList, err := suite.service.readPeopleListFromReader(&buffer)
	assertT.NoError(err)
	assertT.Len(peopleList, 2)
}