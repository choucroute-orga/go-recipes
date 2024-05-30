package api

import (
	"context"
	"net/http"
	"recipes/db"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (api *ApiHandler) getRecipeByID(c echo.Context) error {
	l := logger.WithField("request", "getRecipeByID")
	id := c.Param("id")
	idObject, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		WarnOnError(l, err, "Invalid ID")
		return NewNotFoundError(err)
	}
	recipe, err := db.FindRecipeByID(l, api.mongo, idObject)
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
	if err := c.Validate(recipe); err != nil {
		FailOnError(l, err, "Validation failed")
		return NewBadRequestError(err)
	}
	recipe.ID = db.NewID()
	err := db.SaveRecipe(l, api.mongo, *recipe)
	if err != nil {
		FailOnError(l, err, "Error when trying to save recipe")
		return NewInternalServerError(err)
	}
	return c.JSON(http.StatusCreated, recipe)
}

func (api *ApiHandler) deleteRecipe(c echo.Context) error {
	l := logger.WithField("request", "deleteRecipeByID")
	id := c.Param("id")
	idObject, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		WarnOnError(l, err, "Invalid ID")
		return NewNotFoundError(err)
	}
	err = db.DeleteRecipeByID(l, api.mongo, idObject)
	if err != nil {
		return NewNotFoundError(err)
	}
	return c.NoContent(http.StatusNoContent)
}

func (api *ApiHandler) updateRecipe(c echo.Context) error {
	l := logger.WithField("request", "updateRecipe")
	recipe := new(db.Recipe)
	if err := c.Bind(recipe); err != nil {
		FailOnError(l, err, "Request binding failed")
		return NewBadRequestError(err)
	}
	if err := c.Validate(recipe); err != nil {
		FailOnError(l, err, "Validation failed")
		return NewBadRequestError(err)
	}
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return NewNotFoundError(err)
	}
	recipe.ID = id
	err = db.UpsertOne(l, api.mongo, recipe)
	if err != nil {
		FailOnError(l, err, "Error when trying to save recipe")
		return NewInternalServerError(err)
	}
	return c.JSON(http.StatusCreated, recipe)
}
