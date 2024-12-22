package handler

import (
	"bytes"
	"fmt"
	"github.com/faruqputraaa/go-ticket/config"
	"github.com/faruqputraaa/go-ticket/internal/entity"
	"github.com/golang-jwt/jwt/v5"
	"html/template"
	"net/http"
	"os"

	"github.com/faruqputraaa/go-ticket/internal/http/dto"
	"github.com/faruqputraaa/go-ticket/internal/service"
	"github.com/faruqputraaa/go-ticket/pkg/response"
	"github.com/labstack/echo/v4"
	"gopkg.in/gomail.v2"
)

type OfferHandler struct {
	offerService service.OfferService
	cfg          *config.Config
}

func NewOfferHandler(offerService service.OfferService, cfg *config.Config) OfferHandler {
	return OfferHandler{offerService, cfg}
}

func (h *OfferHandler) sendOfferEmail(offer dto.CreateOfferRequest, subject string, status string, toEmail string, ctx echo.Context) error {

	tmplPath := "templates/email/email_templates.html"
	tmpl, err := os.ReadFile(tmplPath)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Failed to read email template: %v", err)))
	}

	t, err := template.New("offerEmail").Parse(string(tmpl))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Failed to parse email template: %v", err)))
	}

	data := struct {
		NameEvent   string
		Description string
		Status      string
		Email       string
	}{
		NameEvent:   offer.NameEvent,
		Description: offer.Description,
		Status:      status,
		Email:       offer.Email,
	}

	var bodyBuffer bytes.Buffer
	err = t.Execute(&bodyBuffer, data)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Failed to execute email template: %v", err)))
	}

	mail := gomail.NewMessage()
	mail.SetHeader("From", h.cfg.SMTPConfig.Email)
	mail.SetHeader("To", toEmail)
	mail.SetHeader("Subject", subject)
	mail.SetBody("text/html", bodyBuffer.String())

	return h.sendEmail(mail)
}

func (h *OfferHandler) GetOffers(ctx echo.Context) error {
	offers, err := h.offerService.GetAll(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Failed to get all offers: %v", err)))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully showing all tickets", offers))
}

func (h *OfferHandler) GetOffer(ctx echo.Context) error {
	var req dto.GetOfferByIDRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, fmt.Sprintf("Failed to parse request body: %v", err)))
	}
	offer, err := h.offerService.GetByID(ctx.Request().Context(), req.IDOffer)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Failed to get offer: %v", err)))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully showing an offer", offer))
}

func (h *OfferHandler) GetOffersByIDUser(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*entity.JWTCustomClaims)

	IDUser := claims.IDUser

	offers, err := h.offerService.GetByIDUser(ctx.Request().Context(), IDUser)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Failed to get offer: %v", err)))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully showing offers by user ID", offers))
}

func (h *OfferHandler) CreateOffer(ctx echo.Context) error {
	var req dto.CreateOfferRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Failed to parse body: %v", err)))
	}

	err := h.offerService.Create(ctx, req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Failed to create offer: %v", err)))
	}

	subject := fmt.Sprintf("Tawaran untuk membuat event konser yang diajukan oleh email : %s", req.Email)
	if err := h.sendOfferEmail(req, subject, "PENDING", "fafaputra999@gmail.com", nil); err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Failed to send email to admin"))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully created an offer", req))
}

func (h *OfferHandler) updateOfferStatus(ctx echo.Context, status string) error {
	var req dto.UpdateOfferRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, fmt.Sprintf("Failed to parse body: %v", err)))
	}
	req.Status = status
	err := h.offerService.Update(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Failed to update offer: %v", err)))
	}

	offer, err := h.offerService.GetByID(ctx.Request().Context(), req.IDOffer)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Failed to get offer: %v", err)))
	}

	if offer.Email == "" {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Failed to get offer: %v", err)))
	}

	offerDTO := dto.CreateOfferRequest{
		NameEvent:   offer.NameEvent,
		Description: offer.Description,
		Email:       offer.Email,
	}

	subject := fmt.Sprintf("Tawaran untuk membuat event konser yang diajukan oleh email: %s", offer.Email)
	if err := h.sendOfferEmail(offerDTO, subject, status, offer.Email, ctx); err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Failed to send email to user"))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully updated offer status and notified the user", nil))
}

func (h *OfferHandler) ApproveOffer(ctx echo.Context) error {
	return h.updateOfferStatus(ctx, "APPROVED")
}

func (h *OfferHandler) RejectOffer(ctx echo.Context) error {
	return h.updateOfferStatus(ctx, "REJECTED")
}
