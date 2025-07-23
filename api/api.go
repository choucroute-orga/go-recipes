package api

import (
	"recipes/configuration"
	"recipes/db"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type ApiHandler struct {
	dbh  *db.DbHandler
	tracer trace.Tracer
	conf *configuration.Configuration
}

func NewApiHandler(dbh *db.DbHandler, conf *configuration.Configuration) *ApiHandler {
	handler := ApiHandler{
		dbh:  dbh,
		tracer: otel.Tracer(conf.OtelServiceName),
		conf: conf,
	}
	return &handler
}

func (api *ApiHandler) Register(v1 *echo.Group, conf *configuration.Configuration) {

	health := v1.Group("/health")
	health.GET("/alive", api.getAliveStatus)
	health.GET("/live", api.getAliveStatus)
	health.GET("/ready", api.getReadyStatus)

	recipes := v1.Group("/recipe")
	recipes.GET("", api.getRecipes)
	recipes.GET("/:id", api.getRecipeByID)
	recipes.GET("/user/:id", api.getRecipesFromAuthor)
	recipes.GET("/ingredient/:id", api.getRecipeByIngredientID)
	recipes.GET("/title/:title", api.getRecipeByTitle)
	recipes.POST("", api.saveRecipe)
	recipes.PUT("/:id", api.updateRecipe)
	recipes.DELETE("/:id", api.deleteRecipe)
}
