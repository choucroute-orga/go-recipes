package db

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Create a Enum Dish to define the type of dish
type Dish string

// It only takes 3 values
const (
	Starter Dish = "starter"
	Main    Dish = "main"
	Dessert Dish = "dessert"
)

// Metadata is a key value pair to store metadata
type Metadata struct {
	Key   string `json:"key" bson:"key"`
	Value string `json:"value" bson:"value"`
}

// Hashmap is a key value pair to store metadata
type Timer struct {
	Name     string `json:"name" bson:"name" validate:"required"`
	Quantity int    `json:"quantity" bson:"quantity" validate:"required,min=1"`
	Units    string `json:"units" bson:"units" validate:"oneof=seconds minutes hours"`
}

// Reference the ingredient in the catalog MS
type Ingredient struct {
	ID       string  `json:"id" bson:"_id" validate:"omitempty"`
	Quantity float64 `json:"quantity" bson:"quantity" validate:"required,min=0.1"`
	Units    string  `json:"units" bson:"units" validate:"oneof=i is cs tbsp tsp g kg"`
}

type Recipe struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name" validate:"required"`
	Author      string             `json:"author" bson:"author" validate:"required"` // TODO See If w do a MS for that
	Description string             `json:"description" bson:"description" validate:"required"`
	Dish        Dish               `json:"dish" bson:"dish" validate:"oneof=starter main dessert"`
	Servings    int                `json:"servings" bson:"servings" validate:"required,min=1"`
	Metadata    map[string]string  `json:"metadata" bson:"metadata" validate:"omitempty"`
	Timers      []Timer            `json:"timers" bson:"timers" validate:"omitempty,dive,required"`
	Steps       []string           `json:"steps" bson:"steps" validate:"required"`
	Ingredients []Ingredient       `json:"ingredients" bson:"ingredients" validate:"required,dive,required"`
}
