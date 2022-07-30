package main

import (
	"godo/internal/helper/router_builder"
	"godo/internal/repository"
	"net/http"
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/cors"
	"godo/internal/configuration"
	"godo/internal/helper/ilog"
)

func main() {
	logger := ilog.MakeLogger()
	daoLogger := ilog.MakeLoggerWithTag("DAO")
	configLogger := ilog.MakeLoggerWithTag("Config")

	config := configuration.LoadConfig(configLogger)

	//repository.CreateAndPopulateDatabase(logger)
	dao := repository.NewDAO(daoLogger)

	rb := router_builder.New(dao, config)
	router := rb.Init()

	srv := &http.Server{
		Addr:         "0.0.0.0:" + config.ApiPort,
		WriteTimeout: time.Second * 15,
		IdleTimeout:  time.Second * 15,
		Handler:      cors.Default().Handler(router),
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}
