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
	eventRepository := repository.NewEventRepository(db)

	//service
	userService := service.NewUserService(userRepository)
	ticketService := service.NewTicketService(ticketRepository)
	eventService := service.NewEventService(eventRepository)

	//handler
	userHandler := handler.NewUserHandler(userService)
	ticketHandler := handler.NewTicketHandler(ticketService)
	eventHandler := handler.NewEventHandler(eventService)

	return router.PublicRoutes(userHandler, ticketHandler, eventHandler)
}

func BuildPrivateRoute(db *gorm.DB) []route.Route {
	return nil
}
