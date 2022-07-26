package router_builder

import (
	"godo/internal/api/services"
	"godo/internal/helper/ilog"
	"godo/internal/repository"
)

type ServicesCollection struct {
	authService    services.AuthService
	accountService services.AccountService
	projectService services.ProjectService
	storyService   services.StoryService
	taskService    services.TaskService
	userService    services.UserService
}

func newServiceCollection(dao repository.DAO, JWTKey string) ServicesCollection {
	authServiceLogger := ilog.MakeLoggerWithTag("AuthService")
	accountServiceLogger := ilog.MakeLoggerWithTag("AccountService")
	accountQueryLogger := ilog.MakeLoggerWithTag("AccountQuery")
	projectServiceLogger := ilog.MakeLoggerWithTag("ProjectService")
	projectQueryLogger := ilog.MakeLoggerWithTag("ProjectRepo")
	storyServiceLogger := ilog.MakeLoggerWithTag("StoryService")
	storyQueryLogger := ilog.MakeLoggerWithTag("StoryRepo")
	taskServiceLogger := ilog.MakeLoggerWithTag("TaskService")
	taskQueryLogger := ilog.MakeLoggerWithTag("TasQuery")
	userServiceLogger := ilog.MakeLoggerWithTag("UserService")
	userQueryLogger := ilog.MakeLoggerWithTag("UserQuery")

	// Initialize the repositories
	accountQuery := dao.NewAccountQuery(accountQueryLogger)
	projectQuery := dao.NewProjectQuery(projectQueryLogger)
	storyQuery := dao.NewStoryQuery(storyQueryLogger)
	taskQuery := dao.NewTaskQuery(taskQueryLogger)
	userQuery := dao.NewApiUserQuery(userQueryLogger)

	// Initialize the services
	authService := services.NewAuthService(userQuery, []byte(JWTKey), authServiceLogger)
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
