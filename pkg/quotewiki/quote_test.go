package quotewiki

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetQuotesFromCompositeName(t *testing.T) {
	assertT := assert.New(t)

	name := Person{
		FullName: "Albert Einstein",
		ReversedName: "Einstein, Albert",
		Link: "/wiki/Albert_Einstein",
	}

	assertT.Equal("Albert Einstein", name.FullName)
	//err := GetQuotesFromCompositeName(name)
	//assertT.NoError(err)
}

//func TestGetTextFieldByTitleName(t *testing.T) {
//	assertT := assert.New(t)
//
//	name := CompositeName{
//		"Albert Einstein",
//		"Einstein, Albert",
//		"/wiki/Albert_Einstein",
//	}
//
//	err := getTextFieldByTitleName(name.TitleName)
//	assertT.NoError(err)
//}