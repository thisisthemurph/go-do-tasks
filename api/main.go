package main

import (
	"godo/internal/helper/router_builder"
	"godo/internal/repository"
	"net/http"
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"godo/configuration"
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

	router.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	srv := &http.Server{
		Addr:         "0.0.0.0:" + config.ApiPort,
		Handler:      router,
		WriteTimeout: time.Second * 15,
		IdleTimeout:  time.Second * 15,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}
