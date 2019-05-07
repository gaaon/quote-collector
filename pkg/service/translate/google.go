package translate

import (
	"github.com/bregydoc/gtranslate"
	"github.com/gaaon/quote-collector/pkg/constant"
	"github.com/gaaon/quote-collector/pkg/model"
	"github.com/gaaon/quote-collector/pkg/repository"
)

func FindTranslationByGoogle(content string) (string, error) {
	return gtranslate.TranslateWithFromTo(
		content,
		gtranslate.FromTo{
			From: "en",
			To:   "ko",
		},
	)
}

func FindTranslationByGoogleAndSave(content string, entity model.QuoteEntity) (string, interface{}, error) {
	translated, err := FindTranslationByGoogle(content)
	if err != nil {
		return "", nil, err
	}

	_id, err := repository.InsertTranslation(model.TranslationEntity{
		Content:   translated,
		Vendor:    constant.GOOGLE,
		QuoteId:   entity.Id,
		CreatorId: entity.CreatorId,
	})

	return translated, _id, err
}