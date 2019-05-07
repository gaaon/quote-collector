package translate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTranslateByNaver(t *testing.T) {
	assertT := assert.New(t)

	content := "hello, world"
	translated, err := FindTranslationByNaver(content)
	assertT.NoError(err)

	assertT.Contains(translated, "안녕")
	assertT.Contains(translated, "세계")
}
