package repository

import (
	"github.com/gaaon/quote-collector/pkg/model"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const testArticlesXmlName = "testdata/enwikiquote-test-pages-articles.xml"

func TestGetMediaWikiFromReader(t *testing.T) {
	assertT := assert.New(t)

	testArticlesXml, err := os.Open(testArticlesXmlName)
	assertT.NoError(err)
	defer testArticlesXml.Close()

	mediaWiki, err := GetMediaWikiFromReader(testArticlesXml)
	assertT.NoError(err)

	assertT.Len(mediaWiki.Pages, 3)

	mainPage := mediaWiki.Pages[0]
	assertT.Equal("Main Page", mainPage.Title)
}

func TestGetPageMapByFullName(t *testing.T) {
	assertT := assert.New(t)

	person := model.Person{
		FullName: "Albert Einstein",
		ReversedName: "Einstein, Albert",
		Link: "/wiki/Albert_Einstein",
	}

	testArticlesXml, err := os.Open(testArticlesXmlName)
	assertT.NoError(err)
	defer testArticlesXml.Close()

	mediaWiki, err := GetMediaWikiFromReader(testArticlesXml)
	assertT.NoError(err)

	pageMap := GetPersonNamePageMapFromMediaWiki(mediaWiki)
	assertT.Len(pageMap, 3)

	page, exists := pageMap[person.FullName]
	assertT.True(exists)
	assertT.Equal(person.FullName, page.Title)
}

func TestGetQuotesByFullName(t *testing.T) {
	assertT := assert.New(t)

	person := model.Person{
		FullName: "Albert Einstein",
		ReversedName: "Einstein, Albert",
		Link: "/wiki/Albert_Einstein",
	}

	testArticlesXml, err := os.Open(testArticlesXmlName)
	assertT.NoError(err)
	defer testArticlesXml.Close()

	mediaWiki, err := GetMediaWikiFromReader(testArticlesXml)
	assertT.NoError(err)

	pageMap := GetPersonNamePageMapFromMediaWiki(mediaWiki)

	quotes, err := FindQuotesInPageMapByFullName(pageMap, person.FullName)
	assertT.NoError(err)
	assertT.True(len(quotes) > 0)
}
