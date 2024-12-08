package handler

import (
	"net/http"

	"github.com/faruqputraaa/go-ticket/internal/http/dto"
	"github.com/faruqputraaa/go-ticket/internal/service"
	"github.com/faruqputraaa/go-ticket/pkg/response"
	"github.com/labstack/echo/v4"
)

type EventHandler struct {
	eventService service.EventService
}

func NewEventHandler(eventService service.EventService) EventHandler {
	return EventHandler{eventService}
}

func (h *EventHandler) GetEvents(ctx echo.Context) error {
	events, err := h.eventService.GetAll(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("succesfully showing all event", events))
}

func (h *EventHandler) CreateEvent(ctx echo.Context) error {
	var req dto.CreateEventRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	err := h.eventService.Create(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("succesfuly create a event", nil))
}

func (h *EventHandler) GetEvent(ctx echo.Context) error {
	var req dto.GetEventByIDRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	event, err := h.eventService.GetByID(ctx.Request().Context(), req.IDEvent)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("succesfuly showing a event", event))
}

func (h *EventHandler) UpdateEvent(ctx echo.Context) error {
	var req dto.UpdateEventRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}
	err := h.eventService.Update(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("succesfuly update a event", nil))
}

func (h *EventHandler) DeleteEvent(ctx echo.Context) error {
	var req dto.GetEventByIDRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	event, err := h.eventService.GetByID(ctx.Request().Context(), req.IDEvent)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	err = h.eventService.Delete(ctx.Request().Context(), event)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("succesfuly delete a event", nil))

}
