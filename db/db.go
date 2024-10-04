package db

import (
	"context"
	"recipes/configuration"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func New(conf *configuration.Configuration) (*mongo.Client, error) {

	// Database connexion

	loger.Info("Connecting to MongoDB..." + conf.DBURI)
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(conf.DBURI))
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
	return client, nil
}
