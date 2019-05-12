package repository

import (
	"github.com/gaaon/quote-collector/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
)

func InsertPerson(fullName string, koreanName string, link string) (interface{}, error) {
	var person = model.Person{FullName: fullName, KoreanName: koreanName, Link: link}

	res, err := peopleCollection.InsertOne(ctx, person)
	if err != nil {
		return nil, err
	}

	return res.InsertedID, nil
}

func InsertPeopleListIntoDB(peopleList []model.Person) error {
	data := make([]interface{}, 0)

	for _, person := range peopleList {
		data = append(data, person)
	}

	_, err := peopleCollection.InsertMany(ctx, data)

	return err
}

func dropPeopleCollection() error {
	err := peopleCollection.Drop(ctx)
	return err
}

func FindPeopleList() ([]model.Person, error) {
	cur, err := peopleCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var peopleList []model.Person

	for cur.Next(ctx) {
		var person model.Person
		err = cur.Decode(&person)

		peopleList = append(peopleList, person)
	}
	if err = cur.Err(); err != nil {
		return nil, err
	}

	return peopleList, nil
}
