package router

import (
	"github.com/faruqputraaa/go-ticket/internal/http/handler"
	"github.com/faruqputraaa/go-ticket/pkg/route"
)

func PublicRoutes(MovieHandler handler.TicketHandler) []route.Route {
	return []route.Route{}

}

func PrivateRoutes() []route.Route {
	return nil
}
