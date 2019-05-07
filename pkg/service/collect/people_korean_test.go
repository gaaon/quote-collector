package collect

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetKoreanNameFromEnglish(t *testing.T) {
	assertT := assert.New(t)
	name := "Abbas, Mahmoud"

	koreanName, err := GetKoreanNameFromEnglish(name)
	assertT.NoError(err)
	assertT.Equal("마흐무드 압바스", koreanName)
}