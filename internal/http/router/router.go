package router

import (
	"net/http"

	"github.com/faruqputraaa/go-ticket/internal/http/handler"
	"github.com/faruqputraaa/go-ticket/pkg/route"
)

var (
	adminOnly = []string{"ADMIN"}
	allRoles  = []string{"ADMIN", "BUYER"}
)

func PublicRoutes(
	userHandler handler.UserHandler,
	ticketHandler handler.TicketHandler,
	eventHandler handler.EventHandler,
) []route.Route {
	return []route.Route{
		{
			Method:  http.MethodPost,
			Path:    "/login",
			Handler: userHandler.Login,
		},
		{
			Method:  http.MethodPost,
			Path:    "/register",
			Handler: userHandler.Register,
		},
		{
			Method:  http.MethodPost,
			Path:    "/request-reset-password",
			Handler: userHandler.ResetPasswordRequest,
		},
		{
			Method:  http.MethodPost,
			Path:    "/reset-password/:token",
			Handler: userHandler.ResetPassword,
		},
		{
			Method:  http.MethodGet,
			Path:    "/verify-email/:token",
			Handler: userHandler.VerifyEmail,
		},

		{
			Method:  http.MethodGet,
			Path:    "/event",
			Handler: eventHandler.GetEvents,
		},
		{
			Method:  http.MethodPost,
			Path:    "/event",
			Handler: eventHandler.CreateEvent,
		},
		{
			Method:  http.MethodGet,
			Path:    "/event/:id_event",
			Handler: eventHandler.GetEvent,
		},
		{
			Method:  http.MethodPut,
			Path:    "/event/:id_event",
			Handler: eventHandler.UpdateEvent,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/event/:id_event",
			Handler: eventHandler.DeleteEvent,
		},
		{
			Method:  http.MethodGet,
			Path:    "/event/search/name",
			Handler: eventHandler.SearchByName,
		},
		{
			Method:  http.MethodGet,
			Path:    "/event/search/location",
			Handler: eventHandler.SearchByLocation,
		},
	}

}

func PrivateRoutes(
	userHandler handler.UserHandler,
	ticketHandler handler.TicketHandler,
	eventHandler handler.EventHandler,
	offerHandler handler.OfferHandler,
) []route.Route {
	return []route.Route{
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: userHandler.GetUsers,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/:id_user",
			Handler: userHandler.GetUser,
			Roles:   adminOnly,
		},
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: userHandler.CreateUser,
			Roles:   adminOnly,
		},
		{
			Method:  http.MethodPut,
			Path:    "/users/:id_user",
			Handler: userHandler.UpdateUser,
			Roles:   adminOnly,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/users/:id_user",
			Handler: userHandler.DeleteUser,
			Roles:   adminOnly,
		},
		{
			Method:  http.MethodGet,
			Path:    "/offers",
			Handler: offerHandler.GetOffers,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPost,
			Path:    "/offers",
			Handler: offerHandler.CreateOffer,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/offers/:id_offer",
			Handler: offerHandler.GetOffer,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPut,
			Path:    "/offer/approve/:id_offer",
			Handler: offerHandler.ApproveOffer,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPut,
			Path:    "/offer/reject/:id_offer",
			Handler: offerHandler.RejectOffer,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/offer/user",
			Handler: offerHandler.GetOffersByIDUser,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPost,
			Path:    "/ticket",
			Handler: ticketHandler.CreateTicket,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/ticket",
			Handler: ticketHandler.GetTickets,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/ticket/:id_ticket",
			Handler: ticketHandler.GetTicket,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/ticket/:id_event",
			Handler: ticketHandler.GetTicket,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPut,
			Path:    "/ticket/:id_ticket",
			Handler: ticketHandler.UpdateTicket,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/ticket/:id_ticket",
			Handler: ticketHandler.DeleteTicket,
			Roles:   allRoles,
		},
	}
}
