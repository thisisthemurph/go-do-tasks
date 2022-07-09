package main

import (
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
	"github.com/sirupsen/logrus"
)

func main() {
	// Make the loggers
	logger := makeLogger()
	configLogger := makeLoggerWithTag("Config")

	config := config.LoadConfig(configLogger)
	repository.CreateAndPopulateDatabase(logger)

	dao := repository.NewDAO(makeLoggerWithTag("DAO"))
	handlers := buildHandlers(dao, logger)
	router := buildRouter(handlers, logger)

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

func buildHandlers(dao repository.DAO, logger logrus.FieldLogger) handler.IHandler {
	// Make loggers with the specific required parameters
	projectServiceLogger := makeLoggerWithTag("ProjectService")
	storyServiceLogger := makeLoggerWithTag("ProjectService")

	// Instantiate the services
	projectService := services.NewProjectService(dao, projectServiceLogger)
	storyService := services.NewStoryService(dao, storyServiceLogger)

	return handler.MakeHandlers(
		makeLoggerWithTag("Handler"),
		&projectService,
		&storyService,
	)
}

func buildRouter(handlers handler.IHandler, logger logrus.FieldLogger) *mux.Router {
	r := mux.NewRouter()
	router := r.PathPrefix("/api").Subrouter()

	middlewareLogger := makeLoggerWithTag("Middleware")
	gm := middleware.NewGenericMiddleware(middlewareLogger)
	pm := middleware.NewProjectMiddleware(middlewareLogger)

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

func makeLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: true,
		PrettyPrint:      true,
	})

	return logger
}

func makeLoggerWithTag(tag string) logrus.FieldLogger {
	logger := makeLogger()
	return logger.WithFields(logrus.Fields{
		"tag": tag,
	})
}
