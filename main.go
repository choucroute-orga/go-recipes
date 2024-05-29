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
	logger.Info("Cacahuete API Starting...")

	conf := configuration.New()

	mongo, err := db.New(conf)

	if err != nil {
		return
	}

	val := validation.New(conf)
	r := api.New(val)
	v1 := r.Group(conf.ListenRoute)

	if err != nil {
		return
	}
	h := api.NewApiHandler(mongo, conf)

	h.Register(v1, conf)
	r.Logger.Fatal(r.Start(fmt.Sprintf("%v:%v", conf.ListenAddress, conf.ListenPort)))
}
