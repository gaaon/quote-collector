package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const defaultDBName = "quote-collector"

var (
	mongoClient           *mongo.Client
	ctx                   context.Context
	peopleCollection      *mongo.Collection
	quoteCollection       *mongo.Collection
	translationCollection *mongo.Collection
)

func updateCollections(dbName string) {
	peopleCollection = mongoClient.Database(dbName).Collection("people")
	quoteCollection = mongoClient.Database(dbName).Collection("quotes")
	translationCollection = mongoClient.Database(dbName).Collection("translations")
}

func init() {
	var err error

	//ctx, _ = context.WithTimeout(context.Background(), 100 * time.Second)
	ctx = context.Background()
	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	updateCollections(defaultDBName)
}

func SetDBName(newDBName string) {
	updateCollections(newDBName)
}

func RecoverDBName() {
	updateCollections(defaultDBName)
}
