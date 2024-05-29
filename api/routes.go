package api

import (
	"context"
	"net/http"
	"recipes/db"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

var logger = logrus.WithField("context", "api/routes")

func (api *ApiHandler) getAliveStatus(c echo.Context) error {
	l := logger.WithField("request", "getAliveStatus")
	status := NewHealthResponse(LiveStatus)
	if err := c.Bind(status); err != nil {
		FailOnError(l, err, "Response binding failed")
		return NewInternalServerError(err)
	}
	l.WithFields(logrus.Fields{
		"action": "getStatus",
		"status": status,
	}).Debug("Health Status ping")

	return c.JSON(http.StatusOK, &status)
}

func (api *ApiHandler) getReadyStatus(c echo.Context) error {
	l := logger.WithField("request", "getReadyStatus")
	err := api.mongo.Ping(context.Background(), nil)
	if err != nil {
		WarnOnError(l, err, "Unable to ping database to check connection.")
		return c.JSON(http.StatusServiceUnavailable, NewHealthResponse(NotReadyStatus))
	}

	return c.JSON(http.StatusOK, NewHealthResponse(ReadyStatus))
}

func (api *ApiHandler) getRecipeByTitle(c echo.Context) error {
	l := logger.WithField("request", "getRecipeByTitle")
	title := c.Param("title")
	recipe, err := db.FindRecipeByTitle(l, api.mongo, title)
	if err != nil {
		return NewNotFoundError(err)
	}
	return c.JSON(http.StatusOK, recipe)

}

func (api *ApiHandler) saveRecipe(c echo.Context) error {
	l := logger.WithField("request", "saveRecipe")
	recipe := new(db.Recipe)
	if err := c.Bind(recipe); err != nil {
		FailOnError(l, err, "Request binding failed")
		return NewInternalServerError(err)
	}
	recipe.ID = db.NewID()
	err := db.SaveRecipe(l, api.mongo, *recipe)
	if err != nil {
		FailOnError(l, err, "Error when trying to save recipe")
		return NewInternalServerError(err)
	}
	return c.JSON(http.StatusCreated, recipe)
}
