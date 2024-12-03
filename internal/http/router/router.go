package router

import (
	"github.com/faruqputraaa/go-ticket/internal/http/handler"
	"github.com/faruqputraaa/go-ticket/pkg/route"
	"net/http"
)

func PublicRoutes(userHandler handler.UserHandler, ticketHandler handler.TicketHandler) []route.Route {
	return []route.Route{
		{
			Method:  http.MethodPost,
			Path:    "/login",
			Handler: userHandler.Login,
		},
		{
			Method:  http.MethodPost,
			Path:    "/ticket",
			Handler: ticketHandler.CreateTicket,
		},
		{
			Method:  http.MethodGet,
			Path:    "/ticket",
			Handler: ticketHandler.GetTickets,
		},
	}

}

func PrivateRoutes() []route.Route {
	return nil
}
