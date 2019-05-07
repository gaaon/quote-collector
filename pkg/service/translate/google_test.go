package translate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTranslateByGoogle(t *testing.T) {
	assertT := assert.New(t)

	content := "hello, world"
	translated, err := TranslateByGoogle(content)
	assertT.NoError(err)

	assertT.Contains(translated, "안녕")
	assertT.Contains(translated, "세상")
}