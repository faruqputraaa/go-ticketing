package builder

import (
	"github.com/faruqputraaa/go-ticket/config"
	"github.com/faruqputraaa/go-ticket/internal/http/handler"
	"github.com/faruqputraaa/go-ticket/internal/http/router"
	"github.com/faruqputraaa/go-ticket/internal/repository"
	"github.com/faruqputraaa/go-ticket/internal/service"
	"github.com/faruqputraaa/go-ticket/pkg/route"
	"gorm.io/gorm"
)

func BuildPublicRoute(cfg *config.Config, db *gorm.DB) []route.Route {
	//repository
	userRepository := repository.NewUserRepository(db)
	ticketRepository := repository.NewTicketRepository(db)
	eventRepository := repository.NewEventRepository(db)

	//service
	userService := service.NewUserService(userRepository)
	tokenService := service.NewTokenService(cfg.JWTConfig.SecretKey)
	ticketService := service.NewTicketService(ticketRepository)
	eventService := service.NewEventService(eventRepository)

	//handler
	userHandler := handler.NewUserHandler(tokenService, userService)
	ticketHandler := handler.NewTicketHandler(ticketService)
	eventHandler := handler.NewEventHandler(eventService)

	return router.PublicRoutes(userHandler, ticketHandler, eventHandler)
}

func BuildPrivateRoute(cfg *config.Config, db *gorm.DB) []route.Route {
	userRepository := repository.NewUserRepository(db)
	ticketRepository := repository.NewTicketRepository(db)
	eventRepository := repository.NewEventRepository(db)

	//service
	userService := service.NewUserService(userRepository)
	tokenService := service.NewTokenService(cfg.JWTConfig.SecretKey)
	ticketService := service.NewTicketService(ticketRepository)
	eventService := service.NewEventService(eventRepository)

	//handler
	userHandler := handler.NewUserHandler(tokenService, userService)
	ticketHandler := handler.NewTicketHandler(ticketService)
	eventHandler := handler.NewEventHandler(eventService)

	return router.PublicRoutes(userHandler, ticketHandler, eventHandler)
}
