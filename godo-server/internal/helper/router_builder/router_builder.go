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
	sc ServiceCollection
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
	projectHandler := handler.NewProjectsHandler(projectLogger, b.sc.projectService, b.sc.tagService)

	b.Post("/project", projectHandler.CreateProject)
	b.Get("/project", projectHandler.GetAllProjects)
	b.Get("/project/{id:[a-f0-9-]+}", projectHandler.GetProjectById)

	// Status
	b.Put("/project/{id:[a-f0-9-]+}/status", projectHandler.UpdateProjectStatus)

	// Tags
	b.Post("/project/{id:[a-f0-9-]+}/tag", projectHandler.AddTagToProject)
	b.Delete("/project/{projectId:[a-f0-9-]+}/tag/{tagId:[0-9]+}", projectHandler.DeleteProjectTag)
}

func (b *routerBuilder) buildStoryRouter() {
	storyLogger := ilog.MakeLoggerWithTag("StoryHandler")
	storyHandler := handler.NewStoriesHandler(storyLogger, b.sc.storyService, b.sc.projectService)

	b.Get("/story", storyHandler.GetAllStories)
	b.Get("/story", storyHandler.CreateStory)
	b.Get("/story/{id:[a-f0-9-]+}", storyHandler.GetStoryById)
	b.Put("/story/{id:[a-f0-9-]+}", storyHandler.UpdateStory)
	b.Delete("/story/{id:[a-f0-9-]+}", storyHandler.DeleteStory)
}

func (b *routerBuilder) buildTaskRouter() {
	taskLogger := ilog.MakeLoggerWithTag("TaskHandler")
	taskHandler := handler.NewTasksHandler(taskLogger, b.sc.taskService, b.sc.tagService)

	b.Post("/task", taskHandler.CreateTask)
	b.Get("/task", taskHandler.GetAllTasks)
	b.Get("/task/{id:[a-f0-9-]+}", taskHandler.GetTaskById)
	b.Put("/task/{id:[a-f0-9-]+}", taskHandler.UpdateTask)

	// Type and status
	b.Put("/task/{id:[a-f0-9-]+}/type", taskHandler.UpdateTaskStatus)
	b.Put("/task/{id:[a-f0-9-]+}/status", taskHandler.UpdateTaskStatus)

	// Tags
	b.Put("/task/{taskId:[a-f0-9-]+}/tag/{tagId:[0-9]+}", taskHandler.AddTag)
	b.Delete("/task/{taskId:[a-f0-9-]+}/tag/{tagId:[0-9]+}", taskHandler.RemoveTag)
}

type HttpHandlerFunc = func(w http.ResponseWriter, r *http.Request)

func (b *routerBuilder) Get(path string, f HttpHandlerFunc) {
	b.ar.HandleFunc(path, f).Methods(http.MethodGet)
}

func (b *routerBuilder) Post(path string, f HttpHandlerFunc) {
	b.ar.HandleFunc(path, f).Methods(http.MethodPost)
}

func (b *routerBuilder) Put(path string, f HttpHandlerFunc) {
	b.ar.HandleFunc(path, f).Methods(http.MethodPut)
}

func (b *routerBuilder) Delete(path string, f HttpHandlerFunc) {
	b.ar.HandleFunc(path, f).Methods(http.MethodDelete)
}
