package main

import (
	"net/http"
	"time"

	"godo/internal/api/handler"
	"godo/internal/api/middleware"
	"godo/internal/api/services"
	"godo/internal/config"
	"godo/internal/helper/ilog"
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
	servicesCollection := buildServices(dao, config)
	router := buildRouter(servicesCollection, logger)

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

type MiddlewareCollection struct {
	Generic middleware.GenericMiddleware
	Auth    middleware.AuthMiddleware
	Project middleware.ProjectMiddleware
}

func buildServices(dao repository.DAO, config config.Config) ServicesCollection {
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
	authService := services.NewAuthService(userQuery, []byte(config.JWTKey), authServiceLogger)
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

func buildMiddleware(logger ilog.StdLogger, services ServicesCollection) MiddlewareCollection {
	return MiddlewareCollection{
		Generic: middleware.NewGenericMiddleware(logger),
		Auth:    middleware.NewAuthMiddleware(logger, services.authService, services.userService),
		Project: middleware.NewProjectMiddleware(logger),
	}
}

func buildRouter(services ServicesCollection, logger logrus.FieldLogger) *mux.Router {
	middlewareLogger := makeLoggerWithTag("Middleware")
	middlewares := buildMiddleware(middlewareLogger, services)

	r := mux.NewRouter()
	router := r.PathPrefix("/api").Subrouter()
	router.Use(middlewares.Generic.LoggingMiddleware)

	configureUserAuthRouter(router, services)
	configureProjectRouter(router, services, middlewares)

	return r
}

func configureUserAuthRouter(router *mux.Router, services ServicesCollection) {
	userHandlerLogger := makeLoggerWithTag("userHandler")
	userHandler := handler.NewUsersHandler(
		userHandlerLogger,
		services.authService,
		services.userService,
	)

	userAuthRouter := router.Methods(http.MethodPost).Subrouter()
	userAuthRouter.HandleFunc("/auth/login", userHandler.Login)
	userAuthRouter.HandleFunc("/auth/register", userHandler.Register)
}

func configureProjectRouter(router *mux.Router, services ServicesCollection, middlewares MiddlewareCollection) {
	projectLogger := makeLoggerWithTag("ProjectHandler")

	projectHandler := handler.NewProjectsHandler(projectLogger, services.projectService)
	projectGetRouter := router.Methods(http.MethodGet).Subrouter()
	projectGetRouter.HandleFunc("/project", projectHandler.GetAllProjects)
	projectGetRouter.HandleFunc("/project/{id:[a-f0-9-]+}", projectHandler.GetProjectById)
	projectGetRouter.Use(middlewares.Auth.AuthenticateRequestMiddleware)

	projectPostRouter := router.Methods(http.MethodPost).Subrouter()
	projectPostRouter.HandleFunc("/project", projectHandler.CreateProject)
	projectPostRouter.Use(middlewares.Auth.AuthenticateRequestMiddleware)
	projectPostRouter.Use(middlewares.Project.ValidateNewProjectDtoMiddleware)
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
