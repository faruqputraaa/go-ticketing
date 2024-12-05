package handler

import (
	"net/http"

	"github.com/faruqputraaa/go-ticket/internal/http/dto"
	"github.com/faruqputraaa/go-ticket/internal/service"
	"github.com/faruqputraaa/go-ticket/pkg/response"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	tokenService service.TokenService
	userService service.UserService
}

func NewUserHandler(
	tokenService service.TokenService,
	userService service.UserService,
	) UserHandler {
	return UserHandler{tokenService,userService}
}

func (h *UserHandler) Login(ctx echo.Context) error {
	var req dto.UserLoginRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	claims, err := h.userService.Login(ctx.Request().Context(), req.Username, req.Password)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	token, err := h.tokenService.GenerateAccessToken(ctx.Request().Context(), *claims)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("success login", map[string]string{"token": token}))
}

func (h *UserHandler) Register(ctx echo.Context) error {
	var req dto.UserRegisterRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	err := h.userService.Register(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("success register", nil))
}

func (h *UserHandler) GetUsers(ctx echo.Context) error {
	users, err := h.userService.GetAll(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("success showing all user", users))
}

func (h *UserHandler) GetUser(ctx echo.Context) error {
	var req dto.GetUserByRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}
	
	user, err := h.userService.GetByID(ctx.Request().Context(), req.IDUser)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("success showing user", user))
}

func (h *UserHandler) CreateUser(ctx echo.Context) error {
	var req dto.CreateUserRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}
	err := h.userService.Create(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("success create a user", nil))
	}

		func (h *UserHandler) UpdateUser(ctx echo.Context) error {
		var req dto.UpdateUserRequest
		if err := ctx.Bind(&req); err != nil {
			return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
		}
		err := h.userService.Update(ctx.Request().Context(), req)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
		}

		return ctx.JSON(http.StatusOK, response.SuccessResponse("success update a user", nil))
	}

		func (h *UserHandler) DeleteUser(ctx echo.Context) error {
		var req dto.GetUserByRequest
		if err := ctx.Bind(&req); err != nil {
			return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
		}

		user, err :=h.userService.GetByID(ctx.Request().Context(), req.IDUser)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
		}

		err = h.userService.Delete(ctx.Request().Context(), user)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
		}

		return ctx.JSON(http.StatusOK, response.SuccessResponse("success delete a user", nil))
		}