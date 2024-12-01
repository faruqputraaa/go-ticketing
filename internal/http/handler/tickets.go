package handler

import (
	"net/http"

	"github.com/faruqputraaa/go-ticket/internal/http/dto"
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

func (h *TicketHandler) GetTicket(ctx echo.Context) error {
	var req dto.GetTicketByIDRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	ticket, err := h.ticketService.GetByID(ctx.Request().Context(), req.IDTicket)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("succesfuly showing a ticket", ticket))
}

func (h *TicketHandler) GetTicketsByIDEvent(ctx echo.Context) error {
	var req dto.GetTicketByIDEventRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "IDevent is required"))
	}

	// Panggil service untuk mendapatkan tiket berdasarkan IDevent
	tickets, err := h.ticketService.GetByIDEvent(ctx.Request().Context(), req.IDEvent)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	// Kembalikan hasil
	return ctx.JSON(http.StatusOK, response.SuccessResponse("successfully showing tickets by event ID", tickets))
}

func (h *TicketHandler) CreateTicket(ctx echo.Context) error {
	var req dto.CreateTicketRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	err := h.ticketService.Create(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("succesfuly create a ticket", nil))
}

func (h *TicketHandler) UpdateTicket(ctx echo.Context) error {
	var req dto.UpdateTicketRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}
	err := h.ticketService.Update(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("succesfuly update a ticket", nil))
}

func (h *TicketHandler) DeleteTicket(ctx echo.Context) error {
	var req dto.GetTicketByIDRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	ticket, err := h.ticketService.GetByID(ctx.Request().Context(), req.IDTicket)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	err = h.ticketService.Delete(ctx.Request().Context(), ticket)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("succesfuly delete a movie", nil))

}
