package handler

import (
	"fmt"
	"github.com/faruqputraaa/go-ticket/internal/entity"
	"github.com/faruqputraaa/go-ticket/internal/http/dto"
	"github.com/faruqputraaa/go-ticket/internal/service"
	"github.com/faruqputraaa/go-ticket/pkg/response"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gopkg.in/gomail.v2"
	"net/http"
)

type TransactionHandler struct {
	transactionService service.TransactionService
}

func NewTransactionHandler(transactionService service.TransactionService) TransactionHandler {
	return TransactionHandler{transactionService}
}

func (h *TransactionHandler) CreateTransaction(ctx echo.Context) error {
	var req dto.CreateTransactionRequest

	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*entity.JWTCustomClaims)

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid request body"))
	}

	transaction, snapResp, err := h.transactionService.Create(ctx.Request().Context(), req, claims)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	respData := map[string]interface{}{
		"transaction_id": transaction.IDTransaction,
		"payment_url":    snapResp.RedirectURL,
		"total_amount":   transaction.TotalPrice,
		"ticket_price":   transaction.TotalPrice / float64(transaction.QuantityTicket),
		"quantity":       transaction.QuantityTicket,
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Transaction created successfully", respData))
}

func (h *TransactionHandler) GetTransactions(ctx echo.Context) error {
	transactions, err := h.transactionService.GetAll(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("succesfully showing all transactions", transactions))
}

func (h *TransactionHandler) GetTransaction(ctx echo.Context) error {
	var req dto.GetTransactionByIDRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	transaction, err := h.transactionService.GetByID(ctx.Request().Context(), req.IDTransaction)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("succesfuly showing a transaction", transaction))
}

func (h *TransactionHandler) GetTransactionByIDUser(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*entity.JWTCustomClaims)

	IDUser := claims.IDUser

	transactions, err := h.transactionService.GetByIDUser(ctx.Request().Context(), IDUser)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Failed to get offer: %v", err)))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully showing transactions by user ID", transactions))
}

func (h *TransactionHandler) HandleMidtransWebhook(ctx echo.Context) error {
	var notificationPayload map[string]interface{}
	if err := ctx.Bind(&notificationPayload); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid notification payload"))
	}

	orderID, ok := notificationPayload["order_id"].(string)
	if !ok {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid order ID"))
	}

	transactionStatus, ok := notificationPayload["transaction_status"].(string)
	if !ok {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid transaction status"))
	}

	transactionMessage := fmt.Sprintf("Payment notification received: %s", transactionStatus)
	if statusMsg, ok := notificationPayload["status_message"].(string); ok && statusMsg != "" {
		transactionMessage = statusMsg
	}

	// Determine status based on transaction status
	var status string
	switch transactionStatus {
	case "capture", "settlement":
		status = "SUCCESS"
	case "pending":
		status = "PENDING"
	case "deny", "cancel", "expire", "failure":
		status = "FAILED"
	default:
		status = "UNKNOWN"
	}

	updateReq := dto.UpdateTransactionRequest{
		IDTransaction: orderID,
		Status:        status,
	}

	if err := h.transactionService.Update(ctx.Request().Context(), updateReq); err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	if err := h.transactionService.LogTransaction(ctx.Request().Context(), orderID, status, transactionMessage); err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Failed to create transaction log"))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Webhook processed successfully", nil))
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
