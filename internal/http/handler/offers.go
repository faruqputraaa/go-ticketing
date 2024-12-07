package handler

import (
	"fmt"
	"github.com/faruqputraaa/go-ticket/config"
	"net/http"

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

	mail := gomail.NewMessage()
	subject := fmt.Sprintf("Tawaran untuk membuat event konser yang diajukan oleh email : %s", req.Email)

	mail.SetHeader("From", h.cfg.SMTPConfig.Email) // Menggunakan h.cfg.SMTPConfig
	mail.SetHeader("To", "fafaputra999@gmail.com") // Ganti dengan penerima yang sesuai
	mail.SetHeader("Subject", subject)

	body := fmt.Sprintf(`
	==========================
	  Tawaran Event Konser  
	==========================
	
	Nama Event    : %s
	Deskripsi     : %s
	
	
	==========================
	Terima kasih atas perhatian Anda.
	`, req.NameEvent, req.Description)

	mail.SetBody("text/plain", body)

	// Mengirim email menggunakan fungsi yang sudah dibuat
	err = h.sendEmail(mail)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Failed to send email"))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully created an offer", req))
}

func (h *OfferHandler) UpdateOffer(ctx echo.Context) error {
	var req dto.UpdateOfferRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}
	err := h.offerService.Update(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("succesfuly update a offer", nil))
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
