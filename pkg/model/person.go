package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Person struct {
	Id *primitive.ObjectID `bson:"_id,omitempty"`
	FullName string `bson:"fullName,omitempty"`
	ReversedName string `bson:"-"`
	Link string `bson:"link,omitempty"`
	KoreanName string `bson:"koreanName,omitempty"`
}

type PeopleSorts []Person

func (c PeopleSorts) Len() int {
	return len(c)
}

func (c PeopleSorts) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c PeopleSorts) Less(i, j int) bool {
	if c[i].FullName == c[j].FullName {
		return c[i].ReversedName < c[j].ReversedName
	}

	return c[i].FullName < c[j].FullName
}