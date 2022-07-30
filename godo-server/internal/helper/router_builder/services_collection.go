package router_builder

import (
	"godo/internal/api/services"
	"godo/internal/helper/ilog"
	"godo/internal/repository"
)

type ServiceCollection struct {
	authService    services.AuthService
	accountService services.AccountService
	projectService services.ProjectService
	storyService   services.StoryService
	tagService     services.TagService
	taskService    services.TaskService
	userService    services.UserService
}

func newServiceCollection(dao repository.DAO, JWTKey string) ServiceCollection {
	authServiceLogger := ilog.MakeLoggerWithTag("AuthService")
	accountServiceLogger := ilog.MakeLoggerWithTag("AccountService")
	accountQueryLogger := ilog.MakeLoggerWithTag("AccountQuery")
	projectQueryLogger := ilog.MakeLoggerWithTag("ProjectRepo")
	projectServiceLogger := ilog.MakeLoggerWithTag("ProjectService")
	storyQueryLogger := ilog.MakeLoggerWithTag("StoryRepo")
	storyServiceLogger := ilog.MakeLoggerWithTag("StoryService")
	tagQueryLogger := ilog.MakeLoggerWithTag("TagRepo")
	tagServiceLogger := ilog.MakeLoggerWithTag("TagService")
	taskQueryLogger := ilog.MakeLoggerWithTag("TaskQuery")
	taskServiceLogger := ilog.MakeLoggerWithTag("TaskService")
	userQueryLogger := ilog.MakeLoggerWithTag("UserQuery")
	userServiceLogger := ilog.MakeLoggerWithTag("UserService")

	// Initialize the repositories
	accountQuery := dao.NewAccountQuery(accountQueryLogger)
	projectQuery := dao.NewProjectQuery(projectQueryLogger)
	storyQuery := dao.NewStoryQuery(storyQueryLogger)
	tagQuery := dao.NewTagQuery(tagQueryLogger)
	taskQuery := dao.NewTaskQuery(taskQueryLogger)
	userQuery := dao.NewApiUserQuery(userQueryLogger)

	// Initialize the services
	authService := services.NewAuthService(userQuery, []byte(JWTKey), authServiceLogger)
	accountService := services.NewAccountService(accountQuery, accountServiceLogger)
	projectService := services.NewProjectService(projectQuery, projectServiceLogger)
	storyService := services.NewStoryService(storyQuery, storyServiceLogger)
	tagService := services.NewTagService(tagQuery, tagServiceLogger)
	taskService := services.NewTaskService(taskQuery, taskServiceLogger)
	userService := services.NewUserService(userQuery, userServiceLogger)

	return ServiceCollection{
		authService,
		accountService,
		projectService,
		storyService,
		tagService,
		taskService,
		userService,
	}
}
