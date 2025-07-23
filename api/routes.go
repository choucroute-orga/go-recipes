package api

import (
	"context"
	"errors"
	"net/http"
	"recipes/db"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/codes"
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
	err := api.dbh.Client.Ping(context.Background(), nil)
	if err != nil {
		WarnOnError(l, err, "Unable to ping database to check connection.")
		return c.JSON(http.StatusServiceUnavailable, NewHealthResponse(NotReadyStatus))
	}

	return c.JSON(http.StatusOK, NewHealthResponse(ReadyStatus))
}

func (api *ApiHandler) getRecipes(c echo.Context) error {
	l := logger.WithField("request", "getRecipes")
	recipes, err := api.dbh.FindAllRecipes(l)
	if err != nil {
		return NewNotFoundError(err)
	}
	return c.JSON(http.StatusOK, recipes)
}

func (api *ApiHandler) getRecipeByTitle(c echo.Context) error {
	l := logger.WithField("request", "getRecipeByTitle")
	title := c.Param("title")
	recipe, err := api.dbh.FindRecipeByTitle(l, title)
	if err != nil {
		return NewNotFoundError(err)
	}
	return c.JSON(http.StatusOK, recipe)

}

func (api *ApiHandler) getRecipeByID(c echo.Context) error {
	l := logger.WithField("request", "getRecipeByID")
	id := c.Param("id")
	recipe, err := api.dbh.FindRecipeByID(l, id)
	if err != nil {
		return NewNotFoundError(err)
	}
	return c.JSON(http.StatusOK, recipe)
}

func (api *ApiHandler) getRecipeByIngredientID(c echo.Context) error {
	l := logger.WithField("request", "getRecipeByIngredientID")
	id := c.Param("id")

	recipes, err := api.dbh.FindRecipesByIngredientID(l, id)
	if len(*recipes) == 0 {
		err = errors.Join(errors.New("no recipe found for ingredient id"), err)
		l.Error(err)
	}
	if err != nil {
		return NewNotFoundError(err)
	}
	return c.JSON(http.StatusOK, recipes)
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
	if recipe.ID.String() == "" {
		recipe.ID = api.dbh.NewID()
	}
	err := api.dbh.SaveRecipe(l, *recipe)
	if err != nil {
		FailOnError(l, err, "Error when trying to save recipe")
		return NewInternalServerError(err)
	}
	return c.JSON(http.StatusCreated, recipe)
}

func (api *ApiHandler) deleteRecipe(c echo.Context) error {
	l := logger.WithField("request", "deleteRecipeByID")
	id := c.Param("id")
	err := api.dbh.DeleteRecipeByID(l, id)
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
	err = api.dbh.UpsertOne(l, recipe)
	if err != nil {
		FailOnError(l, err, "Error when trying to save recipe")
		return NewInternalServerError(err)
	}
	return c.JSON(http.StatusCreated, recipe)
}

func (api *ApiHandler) getRecipesFromAuthor(c echo.Context) error {
	ctx, span := api.tracer.Start(c.Request().Context(), "api.getRecipesFromAuthor")
	defer span.End()
	l := logger.WithContext(ctx).WithField("request", "getRecipesFromAuthor")

	idParam := new(IDParam)
	if err := c.Bind(idParam); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Request binding failed")
		FailOnError(l, err, "Request binding failed")
		return NewBadRequestError(err)
	}
	if err := c.Validate(idParam); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Request validation failed")
		FailOnError(l, err, "Request validation failed")
		return NewBadRequestError(err)
	}

	recipes, err := api.dbh.FindRecipesByAuthorID(l, idParam.ID)
	if err != nil {
		return NewInternalServerError(err)
	}

	return c.JSON(http.StatusOK, recipes)
}
