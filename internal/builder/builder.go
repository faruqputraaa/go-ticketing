package builder

import (
	"github.com/faruqputraaa/go-ticket/internal/http/handler"
	"github.com/faruqputraaa/go-ticket/internal/http/router"
	"github.com/faruqputraaa/go-ticket/internal/repository"
	"github.com/faruqputraaa/go-ticket/internal/service"
	"github.com/faruqputraaa/go-ticket/pkg/route"
	"gorm.io/gorm"
)

func BuildPublicRoute(db *gorm.DB) []route.Route {
	//repository
	userRepository := repository.NewUserRepository(db)

	//service
	userService := service.NewUserService(userRepository)

	//handler
	userHandler := handler.NewUserHandler(userService)

	return router.PublicRoutes(userHandler)
}

func BuildPrivateRoute(db *gorm.DB) []route.Route {
	return nil
}
