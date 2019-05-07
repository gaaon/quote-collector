package translate

import "github.com/bregydoc/gtranslate"

func TranslateByGoogle(content string) (string, error) {
	return gtranslate.TranslateWithFromTo(
		content,
		gtranslate.FromTo{
			From: "en",
			To:   "ko",
		},
	)
}