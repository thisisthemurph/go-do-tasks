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

	c := config.LoadConfig(configLogger)
	repository.CreateAndPopulateDatabase(logger)

	dao := repository.NewDAO(daoLogger)
	servicesCollection := buildServices(dao, c)
	router := buildRouter(servicesCollection, logger)

	srv := &http.Server{
		Addr:         "0.0.0.0:" + c.ApiPort,
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
	accountService services.AccountService
	projectService services.ProjectService
	storyService   services.StoryService
	taskService    services.TaskService
	userService    services.UserService
}

type MiddlewareCollection struct {
	Generic middleware.GenericMiddleware
	Account middleware.AccountMiddleware
	Auth    middleware.AuthMiddleware
	Project middleware.ProjectMiddleware
}

func buildServices(dao repository.DAO, config config.Config) ServicesCollection {
	// Make loggers with the specific required parameters
	authServiceLogger := makeLoggerWithTag("AuthService")
	accountServiceLogger := makeLoggerWithTag("AccountService")
	accountQueryLogger := makeLoggerWithTag("AccountQuery")
	projectServiceLogger := makeLoggerWithTag("ProjectService")
	projectQueryLogger := makeLoggerWithTag("ProjectRepo")
	storyServiceLogger := makeLoggerWithTag("StoryService")
	storyQueryLogger := makeLoggerWithTag("StoryRepo")
	taskServiceLogger := makeLoggerWithTag("TaskService")
	taskQueryLogger := makeLoggerWithTag("TasQuery")
	userServiceLogger := makeLoggerWithTag("UserService")
	userQueryLogger := makeLoggerWithTag("UserQuery")

	// Initialize the repositories
	accountQuery := dao.NewAccountQuery(accountQueryLogger)
	projectQuery := dao.NewProjectQuery(projectQueryLogger)
	storyQuery := dao.NewStoryQuery(storyQueryLogger)
	taskQuery := dao.NewTaskQuery(taskQueryLogger)
	userQuery := dao.NewApiUserQuery(userQueryLogger)

	// Initialize the services
	authService := services.NewAuthService(userQuery, []byte(config.JWTKey), authServiceLogger)
	accountService := services.NewAccountService(accountQuery, accountServiceLogger)
	projectService := services.NewProjectService(projectQuery, projectServiceLogger)
	storyService := services.NewStoryService(storyQuery, storyServiceLogger)
	taskService := services.NewTaskService(taskQuery, taskServiceLogger)
	userService := services.NewUserService(userQuery, userServiceLogger)

	return ServicesCollection{
		authService,
		accountService,
		projectService,
		storyService,
		taskService,
		userService,
	}
}

func buildMiddleware(logger ilog.StdLogger, services ServicesCollection) MiddlewareCollection {
	return MiddlewareCollection{
		Generic: middleware.NewGenericMiddleware(logger),
		Account: middleware.NewAccountMiddleware(logger),
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

	configureAccountRouter(router, services, middlewares)
	configureUserAuthRouter(router, services)
	configureProjectRouter(router, services, middlewares)
	configureStoryRouter(router, services, middlewares)
	configureTaskRouter(router, services, middlewares)

	return r
}

func configureAccountRouter(router *mux.Router, services ServicesCollection, middlewares MiddlewareCollection) {
	accountHandlerLogger := makeLoggerWithTag("AccountHandler")
	accountHandler := handler.NewAccountsHandler(
		accountHandlerLogger,
		services.accountService,
		services.userService,
	)

	accountPostRouter := router.Methods(http.MethodPost).Subrouter()
	accountPostRouter.HandleFunc("/account", accountHandler.CreateAccount)
	accountPostRouter.Use(middlewares.Account.ValidateNewAccountDtoMiddleware)
}

func configureUserAuthRouter(router *mux.Router, services ServicesCollection) {
	userHandlerLogger := makeLoggerWithTag("UserHandler")
	userHandler := handler.NewUsersHandler(
		userHandlerLogger,
		services.authService,
		services.accountService,
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

	projectPutRouter := router.Methods(http.MethodPut).Subrouter()
	projectPutRouter.HandleFunc("/project/{id:[a-f0-9-]+}/status", projectHandler.UpdateProjectStatus)
	projectPutRouter.Use(middlewares.Auth.AuthenticateRequestMiddleware)
}

func configureStoryRouter(router *mux.Router, services ServicesCollection, middlewares MiddlewareCollection) {
	storyLogger := makeLoggerWithTag("StoryHandler")
	storyHandler := handler.NewStoriesHandler(storyLogger, services.storyService, services.projectService)

	r := router.PathPrefix("/story").Subrouter()
	r.Use(middlewares.Auth.AuthenticateRequestMiddleware)

	r.HandleFunc("", storyHandler.CreateStory).Methods(http.MethodPost)
	r.HandleFunc("", storyHandler.GetAllStories).Methods(http.MethodGet)
	r.HandleFunc("/{id:[a-f0-9-]+}", storyHandler.GetStoryById).Methods(http.MethodGet)
	r.HandleFunc("/{id:[a-f0-9-]+}", storyHandler.UpdateStory).Methods(http.MethodPut)
	r.HandleFunc("/{id:[a-f0-9-]+}", storyHandler.DeleteStory).Methods(http.MethodDelete)

	//storyGetRouter := router.Methods(http.MethodGet).Subrouter()
	//storyGetRouter.HandleFunc("/story", storyHandler.GetAllStories)
	//storyGetRouter.HandleFunc("/story/{id:[a-f0-9-]+}", storyHandler.GetStoryById)
	//storyGetRouter.Use(middlewares.Auth.AuthenticateRequestMiddleware)
	//
	//storyPostRouter := router.Methods(http.MethodPost).Subrouter()
	//storyPostRouter.HandleFunc("/story", storyHandler.CreateStory)
	//storyPostRouter.Use(middlewares.Auth.AuthenticateRequestMiddleware)
	//
	//storyPutRouter := router.Methods(http.MethodPut).Subrouter()
	//storyPutRouter.HandleFunc("/story/{id:[a-f0-9-]+}", storyHandler.UpdateStory)
	//storyPutRouter.Use(middlewares.Auth.AuthenticateRequestMiddleware)
	//
	//storyDeleteRouter := router.Methods(http.MethodDelete).Subrouter()
	//storyDeleteRouter.HandleFunc("story/{id:[a-f0-9-]+}", storyHandler.DeleteStory)
	//storyDeleteRouter.Use(middlewares.Auth.AuthenticateRequestMiddleware)
}

func configureTaskRouter(router *mux.Router, services ServicesCollection, middlewares MiddlewareCollection) {
	taskLogger := makeLoggerWithTag("TaskHandler")
	taskHandler := handler.NewTasksHandler(taskLogger, services.taskService)

	taskGetRouter := router.Methods(http.MethodGet).Subrouter()
	taskGetRouter.HandleFunc("/task", taskHandler.GetAllTasks)
	taskGetRouter.HandleFunc("/task/{id:[a-f0-9-]+}", taskHandler.GetTaskById)
	taskGetRouter.Use(middlewares.Auth.AuthenticateRequestMiddleware)

	taskPostRouter := router.Methods(http.MethodPost).Subrouter()
	taskPostRouter.HandleFunc("/task", taskHandler.CreateTask)
	taskPostRouter.Use(middlewares.Auth.AuthenticateRequestMiddleware)
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
