package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type TranslationEntity struct {
	Id *primitive.ObjectID `bson:"_id,omitempty"`
	QuoteId *primitive.ObjectID `bson:"quoteId,omitempty"`
	Content string `bson:"content"`
	Vendor string `bson:"vendor"`
}