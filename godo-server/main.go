package main

import (
	"log"
	"net/http"
	"time"

	"godo/internal/api/handler"
	"godo/internal/api/middleware"
	"godo/internal/api/services"
	"godo/internal/config"
	"godo/internal/repository"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/cors"
)

func main() {
	config := config.LoadConfig()
	repository.CreateAndPopulateDatabase()

	dao := repository.NewDAO()
	handlers := buildHandlers(dao)
	router := buildRouter(handlers)

	srv := &http.Server{
		Addr:         "0.0.0.0:" + config.ApiPort,
		WriteTimeout: time.Second * 15,
		IdleTimeout:  time.Second * 15,
		Handler:      cors.Default().Handler(router),
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

func buildRouter(handlers handler.IHandler) *mux.Router {
	r := mux.NewRouter()
	router := r.PathPrefix("/api").Subrouter()

	gm := middleware.GenericMiddleware{}
	pm := middleware.ProjectMiddleware{}

	router.Use(gm.LoggingMiddleware)
	router.HandleFunc("/project", handlers.ProjectHandler).Methods(http.MethodGet)
	router.HandleFunc("/project/{id:[a-f0-9-]+}", handlers.ProjectHandler).Methods(http.MethodGet, http.MethodDelete)

	projectWithBodySubRouter := router.Methods(http.MethodPost, http.MethodPut).Subrouter()
	projectWithBodySubRouter.HandleFunc("/project", handlers.ProjectHandler)
	projectWithBodySubRouter.HandleFunc("/project/{id:[a-f0-9-]+}", handlers.ProjectHandler)
	projectWithBodySubRouter.Use(pm.ValidateProjectMiddleware)

	router.HandleFunc("/story", handlers.StoryHandler).Methods(http.MethodGet)
	router.HandleFunc("/story/{id:[a-f0-9-]+}", handlers.StoryHandler)

	return r
}
