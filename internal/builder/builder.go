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
	transactionRepository := repository.NewTransactionRepository(db)

	//service
	userService := service.NewUserService(cfg, userRepository)
	tokenService := service.NewTokenService(cfg.JWTConfig.SecretKey)
	ticketService := service.NewTicketService(ticketRepository)
	eventService := service.NewEventService(eventRepository)
	transactionService := service.NewTransactionService(cfg, transactionRepository)

	//handler
	userHandler := handler.NewUserHandler(tokenService, userService)
	ticketHandler := handler.NewTicketHandler(ticketService)
	eventHandler := handler.NewEventHandler(eventService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	return router.PublicRoutes(userHandler, ticketHandler, eventHandler, transactionHandler)
}

func BuildPrivateRoute(cfg *config.Config, db *gorm.DB) []route.Route {
	//repository
	userRepository := repository.NewUserRepository(db)
	ticketRepository := repository.NewTicketRepository(db)
	eventRepository := repository.NewEventRepository(db)
	offerRepository := repository.NewOfferRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)

	//service
	userService := service.NewUserService(cfg, userRepository)
	tokenService := service.NewTokenService(cfg.JWTConfig.SecretKey)
	ticketService := service.NewTicketService(ticketRepository)
	eventService := service.NewEventService(eventRepository)
	offerService := service.NewOfferService(offerRepository)
	transactionService := service.NewTransactionService(cfg, transactionRepository)

	//handler
	userHandler := handler.NewUserHandler(tokenService, userService)
	ticketHandler := handler.NewTicketHandler(ticketService)
	eventHandler := handler.NewEventHandler(eventService)
	offerHandler := handler.NewOfferHandler(offerService, cfg)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	return router.PrivateRoutes(userHandler, ticketHandler, eventHandler, offerHandler, transactionHandler)
}
