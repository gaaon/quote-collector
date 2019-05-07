package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type QuoteEntity struct {
	Id          *primitive.ObjectID `bson:"_id,omitempty"`
	Content     string              `bson:"content,omitempty"`     // original content
	SubContents []string            `bson:"subContents,omitempty"` // sub contents includes translate, description
	CreatorId   *primitive.ObjectID `bson:"creatorId"`
}

type Quote struct {
	Content     string   // original content
	SubContents []string // sub contents includes translate, description
}

