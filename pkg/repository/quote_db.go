package repository

import (
	"github.com/gaaon/quote-collector/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
)

func InsertQuoteEntitiesIntoDB(quotes []model.QuoteEntity) error {
	data := make([]interface{}, 0)

	for _, quote := range quotes {
		data = append(data, quote)
	}
	_, err := quoteCollection.InsertMany(ctx, data)

	return err
}

func GetQuoteEntitiesWithPerson(quotes []model.Quote, person model.Person) (quoteEntities []model.QuoteEntity) {
	for _, quote := range quotes {
		quoteEntities = append(quoteEntities, model.QuoteEntity{
			Content:     quote.Content,
			SubContents: quote.SubContents,
			CreatorId:   person.Id,
		})
	}

	return
}

func FindQuoteEntitiesFromDB() ([]model.QuoteEntity, error) {
	cur, err := quoteCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var quoteEntities []model.QuoteEntity

	for cur.Next(ctx) {
		var quoteEntity model.QuoteEntity
		err = cur.Decode(&quoteEntity)

		quoteEntities = append(quoteEntities, quoteEntity)
	}
	if err = cur.Err(); err != nil {
		return nil, err
	}

	return quoteEntities, nil
}
