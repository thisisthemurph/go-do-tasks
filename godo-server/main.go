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
	daoLogger := makeLoggerWithTag("DAO")

	config := config.LoadConfig(configLogger)
	repository.CreateAndPopulateDatabase(logger)

	dao := repository.NewDAO(daoLogger)
	servicesCollection := buildServices(dao)
	handlers := buildHandlers(servicesCollection, logger)
	router := buildRouter(servicesCollection, handlers, logger)

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

type ServicesCollection struct {
	authService    services.AuthService
	userService    services.UserService
	projectService services.ProjectService
	storyService   services.StoryService
}

func buildServices(dao repository.DAO) ServicesCollection {
	// Make loggers with the specific required parameters
	authServiceLogger := makeLoggerWithTag("AuthService")
	userServiceLogger := makeLoggerWithTag("UserService")
	userQueryLogger := makeLoggerWithTag("UserQuery")
	projectServiceLogger := makeLoggerWithTag("ProjectService")
	storyServiceLogger := makeLoggerWithTag("ProjectService")
	projectQueryLogger := makeLoggerWithTag("ProjectRepo")
	storyQueryLogger := makeLoggerWithTag("StoryRepo")

	// Initialize the repositories
	userQuery := dao.NewApiUserQuery(userQueryLogger)
	projectQuery := dao.NewProjectQuery(projectQueryLogger)
	storyQuery := dao.NewStoryQuery(storyQueryLogger)

	// Initialize the services
	authService := services.NewAuthService(userQuery, authServiceLogger)
	userService := services.NewUserService(userQuery, userServiceLogger)
	projectService := services.NewProjectService(projectQuery, projectServiceLogger)
	storyService := services.NewStoryService(storyQuery, storyServiceLogger)

	return ServicesCollection{
		authService,
		userService,
		projectService,
		storyService,
	}
}

func buildHandlers(collection ServicesCollection, logger logrus.FieldLogger) handler.IHandler {
	handlerLogger := makeLoggerWithTag("Handler")

	return handler.MakeHandlers(
		handlerLogger,
		&collection.projectService,
		&collection.storyService,
		&collection.authService,
		&collection.userService,
	)
}

func buildRouter(collection ServicesCollection, handlers handler.IHandler, logger logrus.FieldLogger) *mux.Router {
	r := mux.NewRouter()
	router := r.PathPrefix("/api").Subrouter()

	middlewareLogger := makeLoggerWithTag("Middleware")
	am := middleware.NewAuthMiddleware(middlewareLogger, collection.authService)
	gm := middleware.NewGenericMiddleware(middlewareLogger)
	pm := middleware.NewProjectMiddleware(middlewareLogger)

	router.Use(gm.LoggingMiddleware)

	projectGetHandler := router.Methods(http.MethodGet, http.MethodDelete).Subrouter()
	projectGetHandler.HandleFunc("/project", handlers.ProjectHandler)
	projectGetHandler.HandleFunc("/project/{id:[a-f0-9-]+}", handlers.ProjectHandler)
	projectGetHandler.Use(am.AuthenticateRequestMiddleware)

	projectWithBodySubRouter := router.Methods(http.MethodPost, http.MethodPut).Subrouter()
	projectWithBodySubRouter.HandleFunc("/project", handlers.ProjectHandler)
	projectWithBodySubRouter.HandleFunc("/project/{id:[a-f0-9-]+}", handlers.ProjectHandler)
	projectWithBodySubRouter.Use(pm.ValidateProjectMiddleware)

	router.HandleFunc("/story", handlers.StoryHandler).Methods(http.MethodGet)
	router.HandleFunc("/story/{id:[a-f0-9-]+}", handlers.StoryHandler)

	userHandlerLogger := makeLoggerWithTag("userHandler")
	userHandler := handler.NewUsersHandler(
		userHandlerLogger,
		collection.authService,
		collection.userService,
	)

	userAuthRouter := router.Methods(http.MethodPost).Subrouter()
	userAuthRouter.HandleFunc("/auth/login", userHandler.Login)
	userAuthRouter.HandleFunc("/auth/register", userHandler.Register)

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
