package router_builder

import (
	"github.com/gorilla/mux"
	"godo/internal/api/handler"
	"godo/internal/configuration"
	"godo/internal/helper/ilog"
	"godo/internal/repository"
	"net/http"
)

type RouterBuilder interface {
	Init() *mux.Router
}

type routerBuilder struct {
	r  *mux.Router // Unauthenticated router
	ar *mux.Router // Authenticated router
	sc ServicesCollection
	mc MiddlewareCollection
}

func New(dao repository.DAO, config configuration.Config) RouterBuilder {
	sc := newServiceCollection(dao, config.JWTKey)
	mc := newMiddlewareCollection(sc)

	router := mux.NewRouter().PathPrefix("/api").Subrouter()
	authedRouter := router.Name("authedRouter").Subrouter()
	authedRouter.Use(mc.Auth.AuthenticateRequestMiddleware)

	return &routerBuilder{
		r:  router,
		ar: authedRouter,
		sc: sc,
		mc: mc,
	}
}

func (b *routerBuilder) Init() *mux.Router {
	b.buildRouters()
	return b.r
}

func (b *routerBuilder) buildRouters() {
	b.buildAccountRouter()
	b.buildProjectRouter()
	b.buildStoryRouter()
	b.buildTaskRouter()
}

func (b *routerBuilder) buildAccountRouter() {
	userHandlerLogger := ilog.MakeLoggerWithTag("UserHandler")
	userHandler := handler.NewUsersHandler(
		userHandlerLogger,
		b.sc.authService,
		b.sc.accountService,
		b.sc.userService,
	)

	b.r.HandleFunc("/auth/login", userHandler.Login).Methods(http.MethodPost)
	b.r.HandleFunc("/auth/register", userHandler.Register).Methods(http.MethodPost)
}

func (b *routerBuilder) buildProjectRouter() {
	projectLogger := ilog.MakeLoggerWithTag("ProjectHandler")
	projectHandler := handler.NewProjectsHandler(projectLogger, b.sc.projectService)

	b.ar.HandleFunc("/project", projectHandler.CreateProject).Methods(http.MethodPost)
	b.ar.HandleFunc("/project", projectHandler.GetAllProjects).Methods(http.MethodGet)
	b.ar.HandleFunc("/project/{id:[a-f0-9-]+}", projectHandler.GetProjectById).Methods(http.MethodGet)
	b.ar.HandleFunc("/project/{id:[a-f0-9-]+}/status", projectHandler.UpdateProjectStatus).Methods(http.MethodPut)
}

func (b *routerBuilder) buildStoryRouter() {
	storyLogger := ilog.MakeLoggerWithTag("StoryHandler")
	storyHandler := handler.NewStoriesHandler(storyLogger, b.sc.storyService, b.sc.projectService)

	b.ar.HandleFunc("/story", storyHandler.CreateStory).Methods(http.MethodPost)
	b.ar.HandleFunc("/story", storyHandler.GetAllStories).Methods(http.MethodGet)
	b.ar.HandleFunc("/story/{id:[a-f0-9-]+}", storyHandler.GetStoryById).Methods(http.MethodGet)
	b.ar.HandleFunc("/story/{id:[a-f0-9-]+}", storyHandler.UpdateStory).Methods(http.MethodPut)
	b.ar.HandleFunc("/story/{id:[a-f0-9-]+}", storyHandler.DeleteStory).Methods(http.MethodDelete)
}

func (b *routerBuilder) buildTaskRouter() {
	taskLogger := ilog.MakeLoggerWithTag("TaskHandler")
	taskHandler := handler.NewTasksHandler(taskLogger, b.sc.taskService)

	b.ar.HandleFunc("/task", taskHandler.CreateTask).Methods(http.MethodPost)
	b.ar.HandleFunc("/task", taskHandler.GetAllTasks).Methods(http.MethodGet)
	b.ar.HandleFunc("/task/{id:[a-f0-9-]+}", taskHandler.GetTaskById).Methods(http.MethodGet)
}
