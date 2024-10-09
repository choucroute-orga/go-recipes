package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DbHandler struct {
	Client                *mongo.Client
	DBName                string
	RecipesCollectionName string
}

func NewDbHandler(client *mongo.Client, dbName string, recipesCollectionName string) *DbHandler {
	handler := DbHandler{
		Client:                client,
		DBName:                dbName,
		RecipesCollectionName: recipesCollectionName,
	}
	return &handler
}

func New(dbUri string, dbName string, recipesCollectionName string) (*DbHandler, error) {

	// Database connexion

	loger.Info("Connecting to MongoDB..." + dbUri)
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(dbUri))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
	loger.Info("Connected to MongoDB!")
	return NewDbHandler(client, dbName, recipesCollectionName), nil
}
