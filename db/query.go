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

func (dbh *DbHandler) NewID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func (dbh *DbHandler) GetRecipeCollection() *mongo.Collection {
	return dbh.Client.Database(dbh.DBName).Collection(dbh.RecipesCollectionName)
}

func (dbh *DbHandler) FindAllRecipes(l *logrus.Entry) (*[]Recipe, error) {
	recipes := make([]Recipe, 0)
	cursor, err := dbh.GetRecipeCollection().Find(context.Background(), bson.M{})
	if err != nil {
		l.WithError(err).Error("Error when trying to find all recipes")
		return nil, err
	}

	err = cursor.All(context.Background(), &recipes)
	if err != nil {
		l.WithError(err).Error("Error when trying to decode all recipes")
		return nil, err
	}
	return &recipes, nil
}

func (dbh *DbHandler) FindRecipesByIngredientID(l *logrus.Entry, id string) (*[]Recipe, error) {
	cursor, err := dbh.GetRecipeCollection().Find(context.Background(), bson.M{"ingredients._id": id})
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

func (dbh *DbHandler) FindRecipeByTitle(l *logrus.Entry, title string) (*Recipe, error) {
	var recipe Recipe
	// Search if the name is in the title
	err := dbh.GetRecipeCollection().FindOne(context.Background(), bson.M{"name": bson.M{"$regex": title, "$options": "i"}}).Decode(&recipe)
	if err != nil {
		l.WithError(err).Error("Error when trying to find recipe by title")
		return nil, err
	}
	return &recipe, nil
}

func (dbh *DbHandler) FindRecipeByID(l *logrus.Entry, id string) (*Recipe, error) {
	// Convert id to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	filter := map[string]primitive.ObjectID{"_id": objectID}
	var recipe Recipe
	err = dbh.GetRecipeCollection().FindOne(context.Background(), filter).Decode(&recipe)
	if err != nil {
		l.WithError(err).Error("Error when trying to find recipe by id")
		return nil, err
	}
	return &recipe, nil
}

func (dbh *DbHandler) FindRecipesByAuthorID(l *logrus.Entry, author string) (*[]Recipe, error) {
	cursor, err := dbh.GetRecipeCollection().Find(context.Background(), bson.M{"author": author})
	if err != nil {
		l.WithError(err).Error("Error when trying to find recipe by author")
		return nil, err
	}
	recipes := make([]Recipe, 0)
	err = cursor.All(context.Background(), &recipes)
	if err != nil {
		l.WithError(err).Error("Error when trying to decode all recipes")
		return nil, err
	}
	return &recipes, nil
}

// TODO Return the saved recipe
func (dbh *DbHandler) SaveRecipe(l *logrus.Entry, recipe Recipe) error {
	_, err := dbh.GetRecipeCollection().InsertOne(context.Background(), recipe)
	if err != nil {
		l.WithError(err).Error("Error when trying to save recipe")
		return err
	}
	return nil
}

func (dbh *DbHandler) DeleteRecipeByID(l *logrus.Entry, id string) error {
	objectID, _ := primitive.ObjectIDFromHex(id)
	filter := map[string]primitive.ObjectID{"_id": objectID}
	res, err := dbh.GetRecipeCollection().DeleteOne(context.Background(), filter)
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

func (dbh *DbHandler) UpsertOne(l *logrus.Entry, recipe *Recipe) error {
	// Convert id to string
	filter := map[string]primitive.ObjectID{"_id": recipe.ID}
	update := map[string]Recipe{"$set": *recipe}
	res, err := dbh.GetRecipeCollection().UpdateOne(context.Background(), filter, update)
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
