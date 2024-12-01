package router

import (
	"net/http"

	"github.com/faruqputraaa/go-ticket/internal/http/handler"
	"github.com/faruqputraaa/go-ticket/pkg/route"
)

func PublicRoutes(TicketHandler handler.TicketHandler) []route.Route {
	return []route.Route{
		{
			Method:  http.MethodGet,
			Path:    "/tickets",
			Handler: TicketHandler.GetTickets,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tickets/:id",
			Handler: TicketHandler.GetTicket,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tickets/:IDevent",
			Handler: TicketHandler.GetTicket,
		},
		{
			Method:  http.MethodPost,
			Path:    "/tickets",
			Handler: TicketHandler.CreateTicket,
		},
		{
			Method:  http.MethodPut,
			Path:    "/tickets/:id",
			Handler: TicketHandler.UpdateTicket,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/tickets/:id",
			Handler: TicketHandler.DeleteTicket,
		},
	}
}

func PrivateRoutes() []route.Route {
	return nil
}
