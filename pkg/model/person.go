package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

type Person struct {
	Id           *primitive.ObjectID `bson:"_id,omitempty"`
	FullName     string              `bson:"fullName,omitempty"`
	ReversedName string              `bson:"-"`
	Link         string              `bson:"link,omitempty"`
	KoreanName   string              `bson:"koreanName,omitempty"`
	Source       string              `bson:"source,omitempty"`
}

type PeopleSorts []Person

func (c PeopleSorts) Len() int {
	return len(c)
}

func (c PeopleSorts) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c PeopleSorts) Less(i, j int) bool {
	return strings.ToLower(c[i].FullName) < strings.ToLower(c[j].FullName)
}
