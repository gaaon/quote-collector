package repository

import (
	"encoding/xml"
	"github.com/gaaon/quote-collector/pkg/model"
	"io"
	"io/ioutil"
	"strings"
)

type Revision struct {
	XMLName xml.Name `xml:"revision"`
	Text string `xml:"text"`
}
type Page struct {
	XMLName xml.Name `xml:"page"`
	Title string `xml:"title"`
	Revision Revision `xml:"revision"`
}

type MediaWiki struct {
	XMLName xml.Name `xml:"mediawiki"`
	Pages []Page `xml:"page"`
}

func GetMediaWikiFromReader(reader io.Reader) (*MediaWiki, error) {
	var xmlBody MediaWiki
	byteValue, _ := ioutil.ReadAll(reader)
	err := xml.Unmarshal(byteValue, &xmlBody)
	if err != nil {
		return nil, err
	}

	return &xmlBody, nil
}

func GetPersonNamePageMapFromMediaWiki(mediaWiki *MediaWiki) map[string]*Page {
	var personNamePageMap = make(map[string]*Page)

	for i, page := range mediaWiki.Pages {
		personNamePageMap[page.Title] = &mediaWiki.Pages[i]
	}

	return personNamePageMap
}

func filterQuoteContent(content string) string{
	content = strings.ReplaceAll(content, "<br/>", " ")
	content = strings.ReplaceAll(content, "<br>", " ")
	content = strings.ReplaceAll(content, "<BR>", " ")
	return strings.ReplaceAll(content, "<br />", " ")
}

func FindQuotesInPageMapByFullName(
	pageMap map[string]*Page, fullName string) (
	quotes []model.Quote, err error) {
		page, exists := pageMap[fullName]

		if !exists {
			return nil, nil
		} else {
			lines := strings.Split(page.Revision.Text, "\n")

			for i, line := range lines {
				if strings.Contains(line, "==") && (
					strings.Contains(strings.ToLower(line), "see also") ||
					strings.Contains(strings.ToLower(line), "external links")) {
						break
				}

				if len(line) > 2 && line[0] == '*' && line[1] != '*' {
					quoteCandLine := strings.TrimSpace(line[1:])

					if quoteCandLine[0:1] == "[" || quoteCandLine[0:2] == "[[" || quoteCandLine[0:1] == ":" {
						continue
					}

					quote := model.Quote{Content: filterQuoteContent(quoteCandLine)}

					for j := i + 1; j < len(lines); j++ {
						subContentCand := lines[j]

						if len(subContentCand) > 2 && subContentCand[0:2] == "**" && subContentCand[0:3] != "***" {
							quote.SubContents = append(quote.SubContents, strings.TrimSpace(subContentCand[2:]))
						} else {
							break
						}
					}

					quotes = append(quotes, quote)
				}
			}
		}

		return
}