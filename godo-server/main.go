package main

import (
	"log"
	"net/http"
	"time"

	"godo/internal/api/handler"
	"godo/internal/config"
	"godo/internal/repository"
	"godo/internal/services"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/cors"
)

func main() {
	config := config.LoadConfig();
	repository.CreateAndPopulateDatabase()

	dao := repository.NewDAO()
	handlers := buildHandlers(dao)

	router := mux.NewRouter()

	// 3dbb43ff-7ca6-426e-88d5-b466cd0b203e

	router.HandleFunc("/api/project", handlers.ProjectHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/project/{id:[a-f0-9-]+}", handlers.ProjectHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/story", handlers.StoryHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/story/{id:[a-f0-9-]+}", handlers.StoryHandler).Methods(http.MethodGet)

	srv := &http.Server{
		Addr: 			"0.0.0.0:" + config.ApiPort,
		WriteTimeout: 	time.Second * 15,
		IdleTimeout:	time.Second * 15,
		Handler: 		cors.Default().Handler(router),	
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func buildHandlers(dao repository.DAO) handler.IHandler {
	projectService := services.NewProjectService(dao)
	storyService := services.NewStoryService(dao)

	return handler.MakeHandlers(
		&projectService,
		&storyService,
	)
}
