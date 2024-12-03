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
	ticketRepository := repository.NewTicketRepository(db)

	//service
	userService := service.NewUserService(userRepository)
	ticketService := service.NewTicketService(ticketRepository)

	//handler
	userHandler := handler.NewUserHandler(userService)
	ticketHandler := handler.NewTicketHandler(ticketService)

	return router.PublicRoutes(userHandler, ticketHandler)
}

func BuildPrivateRoute(db *gorm.DB) []route.Route {
	return nil
}
