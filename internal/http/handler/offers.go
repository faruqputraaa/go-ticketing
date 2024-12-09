package handler

import (
	"bytes"
	"fmt"
	"github.com/faruqputraaa/go-ticket/config"
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

func (h *OfferHandler) errorResponse(ctx echo.Context, status int, message string) error {
	return ctx.JSON(status, response.ErrorResponse(status, message))
}

func (h *OfferHandler) sendOfferEmail(offer dto.CreateOfferRequest, subject string, status string, toEmail string) error {
	tmplPath := "templates/email_templates.html"
	tmpl, err := os.ReadFile(tmplPath)
	if err != nil {
		return h.errorResponse(nil, http.StatusInternalServerError, fmt.Sprintf("Failed to read email template: %v", err))
	}

	t, err := template.New("offerEmail").Parse(string(tmpl))
	if err != nil {
		return h.errorResponse(nil, http.StatusInternalServerError, fmt.Sprintf("Failed to parse email template: %v", err))
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
		return h.errorResponse(nil, http.StatusInternalServerError, fmt.Sprintf("Failed to execute email template: %v", err))
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
		return h.errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully showing all tickets", offers))
}

func (h *OfferHandler) GetOffer(ctx echo.Context) error {
	var req dto.GetOfferByIDRequest
	if err := ctx.Bind(&req); err != nil {
		return h.errorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	offer, err := h.offerService.GetByID(ctx.Request().Context(), req.IDOffer)
	if err != nil {
		return h.errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully showing an offer", offer))
}

func (h *OfferHandler) GetOffersByIDUser(ctx echo.Context) error {
	var req dto.GetOfferByIDUserRequest
	if err := ctx.Bind(&req); err != nil {
		return h.errorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	offers, err := h.offerService.GetByIDUser(ctx.Request().Context(), req.IDUser)
	if err != nil {
		return h.errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully showing tickets by event ID", offers))
}

func (h *OfferHandler) CreateOffer(ctx echo.Context) error {
	var req dto.CreateOfferRequest
	if err := ctx.Bind(&req); err != nil {
		return h.errorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	err := h.offerService.Create(ctx.Request().Context(), req)
	if err != nil {
		return h.errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	subject := fmt.Sprintf("Tawaran untuk membuat event konser yang diajukan oleh email : %s", req.Email)
	if err := h.sendOfferEmail(req, subject, "PENDING", "fafaputra999@gmail.com"); err != nil {
		return h.errorResponse(ctx, http.StatusInternalServerError, "Failed to send email to admin")
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully created an offer", req))
}

func (h *OfferHandler) updateOfferStatus(ctx echo.Context, status string) error {
	var req dto.UpdateOfferRequest
	if err := ctx.Bind(&req); err != nil {
		return h.errorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	req.Status = status
	err := h.offerService.Update(ctx.Request().Context(), req)
	if err != nil {
		return h.errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	offer, err := h.offerService.GetByID(ctx.Request().Context(), req.IDOffer)
	if err != nil {
		return h.errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	if offer.Email == "" {
		return h.errorResponse(ctx, http.StatusBadRequest, "Invalid email on offer")
	}

	offerDTO := dto.CreateOfferRequest{
		NameEvent:   offer.NameEvent,
		Description: offer.Description,
		Email:       offer.Email,
	}

	subject := fmt.Sprintf("Tawaran untuk membuat event konser yang diajukan oleh email: %s", offer.Email)
	if err := h.sendOfferEmail(offerDTO, subject, status, offer.Email); err != nil {
		return h.errorResponse(ctx, http.StatusInternalServerError, "Failed to send email to user")
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully updated offer status and notified the user", nil))
}

func (h *OfferHandler) ApproveOffer(ctx echo.Context) error {
	return h.updateOfferStatus(ctx, "APPROVED")
}

func (h *OfferHandler) RejectOffer(ctx echo.Context) error {
	return h.updateOfferStatus(ctx, "REJECTED")
}

func (h *OfferHandler) sendEmail(mail *gomail.Message) error {
	dialer := gomail.NewDialer(
		h.cfg.SMTPConfig.Host,
		h.cfg.SMTPConfig.Port,
		h.cfg.SMTPConfig.Email,
		h.cfg.SMTPConfig.Password,
	)
	if err := dialer.DialAndSend(mail); err != nil {
		return err
	}
	return nil
}
