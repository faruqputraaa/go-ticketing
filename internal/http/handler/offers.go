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

func (h *OfferHandler) GetOffers(ctx echo.Context) error {
	offers, err := h.offerService.GetAll(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("succesfully showing all tickets", offers))
}

func (h *OfferHandler) GetOffer(ctx echo.Context) error {
	var req dto.GetOfferByIDRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	offer, err := h.offerService.GetByID(ctx.Request().Context(), req.IDOffer)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("succesfuly showing a offer", offer))
}

func (h *OfferHandler) GetOffersByIDUser(ctx echo.Context) error {
	var req dto.GetOfferByIDUserRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	offers, err := h.offerService.GetByIDUser(ctx.Request().Context(), req.IDUser)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("successfully showing tickets by event ID", offers))
}

func (h *OfferHandler) CreateOffer(ctx echo.Context) error {
	var req dto.CreateOfferRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	err := h.offerService.Create(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	// Update the path to the template file (templates/email_template.html)
	tmplPath := "templates/email_templates.html"
	tmpl, err := os.ReadFile(tmplPath)
	if err != nil {
		// Print the full error to help debug the issue
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Failed to read email template from %s: %v", tmplPath, err)))
	}

	// Parsing the HTML template
	t, err := template.New("offerEmail").Parse(string(tmpl))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Failed to parse email template: %v", err)))
	}

	// Prepare the data to be inserted into the template
	data := struct {
		NameEvent   string
		Description string
		Email       string
	}{
		NameEvent:   req.NameEvent,
		Description: req.Description,
		Email:       req.Email,
	}

	// Create a buffer to store the result of executing the template
	var bodyBuffer bytes.Buffer

	// Execute the template with the provided data
	err = t.Execute(&bodyBuffer, data)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Failed to execute email template: %v", err)))
	}

	// Get the email body as a string from the buffer
	body := bodyBuffer.String()

	// Create the email
	mail := gomail.NewMessage()
	subject := fmt.Sprintf("Tawaran untuk membuat event konser yang diajukan oleh email : %s", req.Email)

	mail.SetHeader("From", h.cfg.SMTPConfig.Email)  // Using SMTPConfig from your config
	mail.SetHeader("To", "gustipadaka19@gmail.com") // Replace with the actual recipient
	mail.SetHeader("Subject", subject)

	mail.SetBody("text/html", body)

	// Send the email
	err = h.sendEmail(mail)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Failed to send email"))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully created an offer", req))
}

func (h *OfferHandler) ApproveOffer(ctx echo.Context) error {
	var req dto.UpdateOfferRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	req.Status = "APPROVED"
	err := h.offerService.Update(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	offer, err := h.offerService.GetByID(ctx.Request().Context(), req.IDOffer)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	if offer.Email == "" {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Email pada offer tidak valid"))
	}

	mail := gomail.NewMessage()
	subject := fmt.Sprintf("Tawaran untuk membuat event konser yang diajukan oleh email: %s", offer.Email)

	fromEmail := h.cfg.SMTPConfig.Email
	if fromEmail == "" {
		fromEmail = "defaultemail@example.com"
	}

	mail.SetHeader("From", fromEmail)
	mail.SetHeader("To", offer.Email)
	mail.SetHeader("Subject", subject)

	body := fmt.Sprintf(`
	==========================
	  Tawaran Event Konser  
	==========================

	Nama Event    : %s
	Deskripsi     : %s
	Status        : %s

	==========================
	Terima kasih atas perhatian Anda.
	`, offer.NameEvent, offer.Description, req.Status)

	mail.SetBody("text/plain", body)

	if err := h.sendEmail(mail); err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Failed to send email: "+err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully updated the offer status to Approved and notified the user", nil))
}

func (h *OfferHandler) RejectOffer(ctx echo.Context) error {
	var req dto.UpdateOfferRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	req.Status = "REJECTED"
	err := h.offerService.Update(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	offer, err := h.offerService.GetByID(ctx.Request().Context(), req.IDOffer)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	if offer.Email == "" {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Email pada offer tidak valid"))
	}

	mail := gomail.NewMessage()
	subject := fmt.Sprintf("Tawaran untuk membuat event konser yang diajukan oleh email: %s", offer.Email)

	fromEmail := h.cfg.SMTPConfig.Email

	mail.SetHeader("From", fromEmail)
	mail.SetHeader("To", offer.Email)
	mail.SetHeader("Subject", subject)

	body := fmt.Sprintf(`
	==========================
	  Tawaran Event Konser  
	==========================

	Nama Event    : %s
	Deskripsi     : %s
	Status        : %s

	==========================
	Terima kasih atas perhatian Anda.
	`, offer.NameEvent, offer.Description, req.Status)

	mail.SetBody("text/plain", body)

	if err := h.sendEmail(mail); err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Failed to send email: "+err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully updated the offer status to Rejected and notified the user", nil))
}

func (h *OfferHandler) sendEmail(mail *gomail.Message) error {
	// Setup koneksi ke SMTP server
	dialer := gomail.NewDialer(
		h.cfg.SMTPConfig.Host,     // SMTP host
		h.cfg.SMTPConfig.Port,     // SMTP port
		h.cfg.SMTPConfig.Email,    // Email pengirim
		h.cfg.SMTPConfig.Password, // Password pengirim
	)

	// Kirim email
	if err := dialer.DialAndSend(mail); err != nil {
		return err
	}
	return nil
}
