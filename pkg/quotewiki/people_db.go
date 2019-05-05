package quotewiki

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const defaultDBName = "quote-collector"
var (
	mongoClient *mongo.Client
	ctx context.Context
	dbName = defaultDBName
	peopleCollection *mongo.Collection
)

type Person struct {
	Id *primitive.ObjectID `bson:"_id,omitempty"`
	FullName string `bson:"fullName,omitempty"`
	ReversedName string `bson:"-"`
	Link string `bson:"link,omitempty"`
	KoreanName string `bson:"koreanName,omitempty"`
}

func SetDBName(newDBName string) {
	dbName = newDBName
	peopleCollection = mongoClient.Database(dbName).Collection("people")
}

func RecoverDBName() {
	dbName = defaultDBName
	peopleCollection = mongoClient.Database(dbName).Collection("people")
}

func init() {
	var err error

	ctx, _ = context.WithTimeout(context.Background(), 10 * time.Second)
	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	peopleCollection = mongoClient.Database(dbName).Collection("people")
}

func InsertPersonIntoDB(fullName string, koreanName string, link string) (interface{}, error) {
	var person = Person{FullName: fullName, KoreanName: koreanName, Link: link}

	res, err := peopleCollection.InsertOne(ctx, person)
	if err != nil {
		return nil, err
	}

	return res.InsertedID, nil
}

func dropPeopleCollection() error {
	err := peopleCollection.Drop(ctx)
	return err
}

func FindPeopleListFromDB() ([]Person, error) {
	cur, err := peopleCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var peopleList []Person

	for cur.Next(ctx) {
		var person Person
		err = cur.Decode(&person)

		peopleList = append(peopleList, person)
	}
	if err = cur.Err(); err != nil {
		return nil, err
	}

	return peopleList, nil
}