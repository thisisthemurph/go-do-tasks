package router_builder

import (
	"godo/internal/api/middleware"
	"godo/internal/helper/ilog"
)

type MiddlewareCollection struct {
	Generic middleware.GenericMiddleware
	Account middleware.AccountMiddleware
	Auth    middleware.AuthMiddleware
	Project middleware.ProjectMiddleware
}

func newMiddlewareCollection(sc ServiceCollection) MiddlewareCollection {
	logger := ilog.MakeLoggerWithTag("Middleware")

	return MiddlewareCollection{
		Generic: middleware.NewGenericMiddleware(logger),
		Account: middleware.NewAccountMiddleware(logger),
		Auth:    middleware.NewAuthMiddleware(logger, sc.authService, sc.userService),
		Project: middleware.NewProjectMiddleware(logger),
	}
}
