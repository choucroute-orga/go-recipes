package db

import (
	"context"
	"errors"

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

func FindAllRecipes(l *logrus.Entry, mongo *mongo.Client) (*[]Recipe, error) {
	collection := mongo.Database("recipe").Collection("recipe")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		l.WithError(err).Error("Error when trying to find all recipes")
		return nil, err
	}
	var recipes []Recipe
	err = cursor.All(context.Background(), &recipes)
	if err != nil {
		l.WithError(err).Error("Error when trying to decode all recipes")
		return nil, err
	}
	return &recipes, nil
}

func FindRecipeByIngredientID(l *logrus.Entry, mongo *mongo.Client, id string) (*[]Recipe, error) {
	collection := mongo.Database("recipe").Collection("recipe")
	cursor, err := collection.Find(context.Background(), bson.M{"ingredients._id": id})
	if err != nil {
		l.WithError(err).Error("Error when trying to find recipe by ingredient id")
		return nil, err
	}
	var recipes []Recipe
	err = cursor.All(context.Background(), &recipes)
	if err != nil {
		l.WithError(err).Error("Error when trying to decode all recipes")
		return nil, err
	}
	return &recipes, nil
}

func FindRecipeByTitle(l *logrus.Entry, mongo *mongo.Client, title string) (*Recipe, error) {
	collection := mongo.Database("recipe").Collection("recipe")
	var recipe Recipe
	// Search if the name is in the title
	err := collection.FindOne(context.Background(), bson.M{"name": bson.M{"$regex": title, "$options": "i"}}).Decode(&recipe)
	if err != nil {
		l.WithError(err).Error("Error when trying to find recipe by title")
		return nil, err
	}
	return &recipe, nil
}

func FindRecipeByID(l *logrus.Entry, mongo *mongo.Client, id primitive.ObjectID) (*Recipe, error) {
	collection := mongo.Database("recipe").Collection("recipe")
	var recipe Recipe
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&recipe)
	if err != nil {
		l.WithError(err).Error("Error when trying to find recipe by id")
		return nil, err
	}
	return &recipe, nil
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

func DeleteRecipeByID(l *logrus.Entry, mongo *mongo.Client, id primitive.ObjectID) error {
	collection := mongo.Database("recipe").Collection("recipe")
	res, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		l.WithError(err).Error("Error when trying to delete recipe by id")
		return err
	}

	if res.DeletedCount == 0 {
		err = errors.New("ID not found")
		l.WithError(err).Error("Error when trying to delete recipe by id")
		return err
	}
	return nil
}

func UpsertOne(l *logrus.Entry, mongo *mongo.Client, recipe *Recipe) error {
	collection := mongo.Database("recipe").Collection("recipe")
	filter := map[string]primitive.ObjectID{"_id": recipe.ID}
	update := map[string]Recipe{"$set": *recipe}
	res, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		l.WithError(err).Error("Error when trying to upsert recipe")
		return err
	}
	if res.MatchedCount == 0 {
		err = errors.New("ID not found")
		l.WithError(err).Error("Error when trying to upsert recipe")
		return err
	}
	return nil
}
