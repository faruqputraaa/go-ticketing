package handler

import (
	"github.com/faruqputraaa/go-ticket/internal/http/dto"
	"github.com/faruqputraaa/go-ticket/internal/service"
	"github.com/faruqputraaa/go-ticket/pkg/response"
	"github.com/labstack/echo/v4"
	"net/http"
)

type EventHandler struct {
	eventHandler service.EventService
}

func NewEventHandler(eventService service.EventService) EventHandler {
	return EventHandler{eventService}
}

func (h *EventHandler) GetEvents(ctx echo.Context) error {
	events, err := h.eventHandler.GetAll(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("succeeffuly showing all events", events))
}

func (h *EventHandler) CreateEvent(ctx echo.Context) error {
	var req dto.CreateEventRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	err := h.eventHandler.Create(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("succeeffuly create event", req))

}
