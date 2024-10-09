package main

import (
	"fmt"
	"recipes/api"
	"recipes/configuration"
	"recipes/db"
	"recipes/validation"

	"github.com/sirupsen/logrus"
)

var logger = logrus.WithFields(logrus.Fields{
	"context": "main",
})

func main() {
	logger.Info("Recipe API Starting...")

	conf := configuration.New()
	logger.Logger.SetLevel(conf.LogLevel)
	dbh, err := db.New(conf.DBURI, conf.DBName, conf.RecipesCollectionName)

	if err != nil {
		return
	}

	val := validation.New(conf)
	r := api.New(val)
	v1 := r.Group(conf.ListenRoute)

	if err != nil {
		return
	}
	h := api.NewApiHandler(dbh, conf)

	h.Register(v1, conf)
	r.Logger.Fatal(r.Start(fmt.Sprintf("%v:%v", conf.ListenAddress, conf.ListenPort)))
}
