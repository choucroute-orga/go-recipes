package api

import (
	"recipes/configuration"
	"recipes/db"

	"github.com/labstack/echo/v4"
)

type ApiHandler struct {
	dbh  *db.DbHandler
	conf *configuration.Configuration
}

func NewApiHandler(dbh *db.DbHandler, conf *configuration.Configuration) *ApiHandler {
	handler := ApiHandler{
		dbh:  dbh,
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
	recipes.GET("/ingredient/:id", api.getRecipeByIngredientID)
	recipes.GET("/title/:title", api.getRecipeByTitle)
	recipes.POST("", api.saveRecipe)
	recipes.PUT("/:id", api.updateRecipe)
	recipes.DELETE("/:id", api.deleteRecipe)
}
