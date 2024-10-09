package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"testing"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func doubleMe(x float64) float64 {
	return x * 2
}

// You can use testing.T, if you want to test the code without benchmarking
func setupSuite(tb testing.TB) func(tb testing.TB) {
	log.Println("setup suite")

	// Return a function to teardown the test
	return func(tb testing.TB) {
		log.Println("teardown suite")
	}
}

// Almost the same as the above, but this one is for single test instead of collection of tests
func setupTest(tb testing.TB) (*DbHandler, func(tb testing.TB)) {
	// log.Println("setup test")

	// return func(tb testing.TB) {
	// 	log.Println("teardown test")
	// }

	// Get a random port for the test, between 1024 and 65535
	exposedPort := fmt.Sprint(rand.Intn(65525-1024) + 1024)
	dbh, pool, resource := InitTestDocker(exposedPort)
	SeedDatabase(dbh.Client)
	return dbh, func(tb testing.TB) {
		CloseTestDocker(dbh.Client, pool, resource)
	}
}

func TestDoubleMe(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	t.Run("Insert one Recipe in the DB", func(t *testing.T) {
		l := logrus.WithField("test", "Insert one Recipe in the DB")
		dbh, teardownTest := setupTest(t)

		recipeJSON := `{
				"id": "59b40d78cc5d6a001237265e",
				"name": "Pate tomates basilic",
				"author": "Arsène Fougerouse",
				"description": "Le lundi c'est spaghetti, le mardi c'est spaghetti... On en mangerait presque toute la semaine !",
				"servings": 4,
				"dish": "main",
				"metadata": {
						"cook time": "30"
				},
				"timers": [
						{
								"name": "preparation time",
								"quantity": 3,
								"units": "minutes"
						},
						{
								"name": "cooking time",
								"quantity": 10,
								"units": "minutes"
						}
				],
				"ingredients": [
						{
								"id": "59b40d78cc5d6a001237265e",
								"quantity": 480,
								"units": "g"
						},
						{
								"id": "5a60f0f6327fe00014912629",
								"quantity": 400,
								"units": "g"
						},
						{
								"id": "598b651ffd078b0011140a21",
								"quantity": 60,
								"units": "g"
						},
						{
								"id": "598b5ebefd078b0011140a17",
								"quantity": 1,
								"units": "i"
						},
						{
								"id": "598b5e26fd078b0011140a16",
								"quantity": 1,
								"units": "i"
						}
				],
				"steps": [
						"Cuire les pâtes en suivant les instructions de préparation du paquet.",
						"Lavez les tomates, puis ajoutez-les dans une poêle à feu moyen avec un filet d'huile d'olive.",
						"Râpez ou émincez l'ail finement et ajoutez-le dans la poêle avec les tomates. Faites revenir les tomates 2 à 3 minutes.",
						"Une fois les tomates cuites, écrasez-les gentiment à l'aide de votre spatule.",
						"Égouttez les pâtes en fin de cuisson puis ajoutez-les dans la poêle avec les tomates.",
						"Râpez le parmesan, ajoutez le basilic émincé, sel, poivre et mélangez. Servir avec un filet d'huile d'olive et quelques feuilles de basilic pour la déco, c'est prêt !"
				]
		}`
		var Recipe1 Recipe
		err := json.Unmarshal([]byte(recipeJSON), &Recipe1)

		if err != nil {
			t.Errorf("Error when trying to unmarshal recipe: %v", err)
		}
		err = dbh.SaveRecipe(l, Recipe1)

		if err != nil {
			t.Errorf("Error when trying to save recipe: %v", err)
		}

		collection := dbh.Client.Database("recipe").Collection("recipe")
		nb, err := collection.CountDocuments(context.Background(), bson.M{})
		if err != nil {
			t.Errorf("Error when trying to find all recipes: %v", err)
		}
		if nb != 1 {
			t.Errorf("Expected 1 recipes, got %v", nb)
		}

		r, err := dbh.FindRecipeByTitle(l, "Pate tomates basilic")

		if err != nil {
			t.Errorf("Error when trying to find recipe by title: %v", err)
		}
		// Assert that the recipe is the same as the one we inserted
		if r.Name != "Pate tomates basilic" {
			t.Errorf("Expected recipe name to be 'Pate tomates basilic', got %v", r.Name)
		}

		defer teardownTest(t)
	})

}
