package handler

import (
	"github.com/faruqputraaa/go-ticket/internal/service"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return UserHandler{userService}
}

func (h *UserHandler) Login(ctx echo.Context) error {
	username := ctx.FormValue("username")
	password := ctx.FormValue("password")

	user, err := h.userService.Login(ctx.Request().Context(), username, password)
	if err != nil {
		return err
	}

	return ctx.JSON(200, user)
}
