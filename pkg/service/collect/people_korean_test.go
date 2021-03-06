package collect

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type NameTranslateServiceTestSuite struct {
	suite.Suite
	service *NameTranslateService
}

func (suite *NameTranslateServiceTestSuite) SetupTest() {
	suite.service = NewNameTranslateService(&http.Client{})
}

func TestNameTranslateServiceTestSuite(t *testing.T) {
	suite.Run(t, new(NameTranslateServiceTestSuite))
}

func (suite *NameTranslateServiceTestSuite) TestGetKoreanNameFromEnglish() {
	assertT := assert.New(suite.T())
	name := "Abbas, Mahmoud"

	koreanName, err := suite.service.TranslateFullNameToKorean(name)
	assertT.NoError(err)
	assertT.Equal("마흐무드 압바스", koreanName)
}
