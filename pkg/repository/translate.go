package repository

import "github.com/gaaon/quote-collector/pkg/model"

func InsertTranslation(translation model.TranslationEntity) (interface{}, error) {
	return translationCollection.InsertOne(ctx, translation)
}
