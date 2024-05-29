package db

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var loger = logrus.WithFields(logrus.Fields{
	"context": "db/query",
})

// func LogAndReturnError(l *logrus.Entry, result *gorm.DB, action string, modelType string) error {
// 	if err := result.Error; err != nil {
// 		l.WithError(err).Error("Error when trying to query database to " + action + " " + modelType)
// 		return err
// 	}
// 	return nil
// }

func NewID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func FindRecipeByTitle(l *logrus.Entry, mongo *mongo.Client, title string) (Recipe, error) {
	collection := mongo.Database("recipe").Collection("recipe")
	var recipe Recipe
	err := collection.FindOne(context.Background(), bson.M{"name": title}).Decode(&recipe)
	if err != nil {
		l.WithError(err).Error("Error when trying to find recipe by title")
		return Recipe{}, err
	}
	return recipe, nil
}

func SaveRecipe(l *logrus.Entry, mongo *mongo.Client, recipe Recipe) error {
	collection := mongo.Database("recipe").Collection("recipe")
	_, err := collection.InsertOne(context.Background(), recipe)
	if err != nil {
		l.WithError(err).Error("Error when trying to save recipe")
		return err
	}
	return nil
}
