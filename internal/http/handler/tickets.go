package handler

import (
	"net/http"

	"github.com/faruqputraaa/go-ticket/internal/service"
	"github.com/faruqputraaa/go-ticket/pkg/response"
	"github.com/labstack/echo/v4"
)

type TicketHandler struct {
	ticketService service.TicketService
}

func NewTicketHandler(ticketService service.TicketService) TicketHandler {
	return TicketHandler{ticketService}
}

func (h *TicketHandler) GetTickets(ctx echo.Context) error {
	tickets, err := h.ticketService.GetAll(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("succesfully showing all tickets", tickets))
}
