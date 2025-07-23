package api

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/codes"
)

type queryDB func(*logrus.Entry, any) (any, error)

type simpleRequest struct {
	Context *echo.Context
	Method  string // Used for tracing, indicate, the name of function that made the call
	Request any    // Should be a pointer to a struct
	Dbfunc  queryDB
}

func (api *ApiHandler) executeSimpleRequest(s *simpleRequest) (any, error) {

	c := s.Context

	// First check with reflect if the request is nil
	if reflect.ValueOf(s.Request).IsNil() {
		return nil, NewBadRequestError(errors.New("Request is nil"))
	}

	ctx, span := api.tracer.Start((*c).Request().Context(), "api."+s.Method)
	defer span.End()
	l := logger.WithContext(ctx).WithField("request", s.Method)

	if s.Request != nil {
		l = l.WithField("requestObject", s.Request)
		// Debug the type and the value of the request
		l.Debug("Trying to bind and validate the Request")
		if err := (*c).Bind(s.Request); err != nil {
			return nil, NewBadRequestError(err)
		}
		if err := (*c).Validate(s.Request); err != nil {
			return nil, NewUnprocessableEntityError(err)
		}
	}

	reqCtx, reqSpan := api.tracer.Start(ctx, fmt.Sprintf("%v.%v", s.Method, "executeDBQuery"))
	l = l.WithContext(reqCtx)
	resp, err := s.Dbfunc(l, s.Request)

	// Check with reflect that the response is a pointer to a struct
	if reflect.ValueOf(resp).IsNil() {
		return nil, NewInternalServerError(errors.New("Response is nil"))
	}

	// Ensure that the response is a pointer to a struct
	if reflect.TypeOf(resp).Kind() != reflect.Ptr {
		return nil, NewInternalServerError(errors.New("Response is not a pointer"))
	}

	// Ensure that the response is a struct
	if reflect.TypeOf(resp).Elem().Kind() != reflect.Struct {
		return nil, NewInternalServerError(errors.New("Response is not a struct"))
	}

	// Convert the response to a struct
	response := reflect.New(reflect.TypeOf(resp).Elem()).Interface()
	reflect.ValueOf(response).Elem().Set(reflect.ValueOf(resp).Elem())

	if resp != nil {
		l = l.WithFields(logrus.Fields{
			"returnObject": resp,
		})
		// reqSpan.SetAttributes(
		// 	attribute.String("returnObject"),
		// )
	}

	if err != nil {
		errMsg := fmt.Sprintf("Error when trying to call method %v with object %v", s.Method, s.Request)
		reqSpan.RecordError(err)
		reqSpan.SetStatus(codes.Error, errMsg)
		FailOnError(l, err, errMsg)
		reqSpan.End()
		return nil, NewInternalServerError(err)
	}
	reqSpan.End()
	l = l.WithContext(ctx)

	l.WithFields(logrus.Fields{
		"responseValue": reflect.ValueOf(response),
		"responseType":  reflect.TypeOf(response),
		"responseKind":  reflect.TypeOf(response).Kind(),
	}).Debug(
		fmt.Sprintf("Response object received from %v function", s.Method),
	)

	return resp, nil
}
