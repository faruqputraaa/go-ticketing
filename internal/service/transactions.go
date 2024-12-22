package service

import (
	"bytes"
	"context"
	_ "errors"
	"fmt"
	"github.com/faruqputraaa/go-ticket/config"
	"github.com/faruqputraaa/go-ticket/internal/entity"
	"github.com/faruqputraaa/go-ticket/internal/http/dto"
	"github.com/faruqputraaa/go-ticket/internal/repository"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"gopkg.in/gomail.v2"
	"text/template"
	"time"
)

type TransactionService interface {
	GetAll(ctx context.Context) ([]entity.Transaction, error)
	GetByID(ctx context.Context, id string) (*entity.Transaction, error)
	GetByIDUser(ctx context.Context, IDUser int) ([]entity.Transaction, error)
	Create(ctx context.Context, req dto.CreateTransactionRequest, claims *entity.JWTCustomClaims) (*entity.Transaction, *snap.Response, error)
	Update(ctx context.Context, req dto.UpdateTransactionRequest) error
	LogTransaction(ctx context.Context, transactionID string, status string, message string) error
	SendSuccessEmail(transactionID string) error
}

type transactionService struct {
	cfg                   *config.Config
	transactionRepository repository.TransactionRepository
}

func NewTransactionService(cfg *config.Config, transactionRepository repository.TransactionRepository) TransactionService {
	return &transactionService{cfg, transactionRepository}
}

func (s *transactionService) Create(ctx context.Context, req dto.CreateTransactionRequest, claims *entity.JWTCustomClaims) (*entity.Transaction, *snap.Response, error) {
	if req.QuantityTicket <= 0 {
		return nil, nil, fmt.Errorf("Quantity of tickets must be greater than 0")
	}

	if req.IDTicket == 0 {
		return nil, nil, fmt.Errorf("ID ticket is required")
	}

	// Fetch ticket and user details
	ticket, err := s.transactionRepository.GetTicketByID(ctx, req.IDTicket)
	if err != nil {
		return nil, nil, fmt.Errorf("Ticket not found")
	}

	_, err = s.transactionRepository.GetUserByID(ctx, req.IDUser)
	if err != nil {
		return nil, nil, fmt.Errorf("User not found")
	}

	var transaction *entity.Transaction
	var snapResp *snap.Response

	// If ticket price is 0, create transaction directly without Midtrans
	if ticket.Price == 0 {
		amount := 0.0
		transactionID := fmt.Sprintf("TRX-%d", time.Now().Unix())

		transaction = &entity.Transaction{
			IDTransaction:  transactionID,
			IDUser:         claims.IDUser,
			IDTicket:       req.IDTicket,
			QuantityTicket: req.QuantityTicket,
			TotalPrice:     amount,
			Status:         "SUCCESS", // Direct success because ticket price is 0
			DateOrder:      time.Now(),
		}

		// Save transaction to database
		if err := s.transactionRepository.Create(ctx, transaction); err != nil {
			return nil, nil, fmt.Errorf("Failed to save transaction")
		}

		// Send success email
		if err := s.SendSuccessEmail(transaction.IDTransaction); err != nil {
			return nil, nil, fmt.Errorf("Failed to send success email: %v", err)
		}
	} else {
		// If ticket price > 0, process the transaction through Midtrans
		amount := float64(ticket.Price) * float64(req.QuantityTicket)
		transactionID := fmt.Sprintf("TRX-%d", time.Now().Unix())

		transaction = &entity.Transaction{
			IDTransaction:  transactionID,
			IDUser:         claims.IDUser,
			IDTicket:       req.IDTicket,
			QuantityTicket: req.QuantityTicket,
			TotalPrice:     amount,
			Status:         "PENDING",
			DateOrder:      time.Now(),
		}

		// Save transaction to database
		if err := s.transactionRepository.Create(ctx, transaction); err != nil {
			return nil, nil, fmt.Errorf("Failed to save transaction")
		}

		// Call Midtrans to create a payment link
		sn := snap.Client{}
		sn.New(s.cfg.MidtransConfig.Serverkey, midtrans.Sandbox)

		reqMidtrans := &snap.Request{
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  transactionID,
				GrossAmt: int64(amount),
			},
		}

		snapResp, err = sn.CreateTransaction(reqMidtrans)
	}

	return transaction, snapResp, nil
}

// GetAll implements TicketService.
func (s *transactionService) GetAll(ctx context.Context) ([]entity.Transaction, error) {
	return s.transactionRepository.GetAll(ctx)
}

// GetByID implements TicketService.
func (s *transactionService) GetByID(ctx context.Context, id string) (*entity.Transaction, error) {
	return s.transactionRepository.GetByID(ctx, id)
}

func (s *transactionService) GetByIDUser(ctx context.Context, IDUser int) ([]entity.Transaction, error) {
	return s.transactionRepository.GetByIdUser(ctx, IDUser)
}

func (s *transactionService) LogTransaction(ctx context.Context, transactionID string, status string, message string) error {
	transactionLog := &entity.TransactionLog{
		IDTransaction: transactionID,
		Status:        status,
		Message:       message,
		CreatedAt:     time.Now(),
	}

	// Call the repository to insert the log into the database
	return s.transactionRepository.CreateLogTransaction(ctx, transactionLog)
}

func (s *transactionService) Update(ctx context.Context, req dto.UpdateTransactionRequest) error {
	transaction, err := s.transactionRepository.GetByID(ctx, req.IDTransaction)
	if err != nil {
		return fmt.Errorf("transaction not found: %v", err)
	}

	transaction.Status = req.Status
	return s.transactionRepository.Update(ctx, transaction)
}

func (s *transactionService) SendSuccessEmail(transactionID string) error {
	// Fetch transaction and user details
	transaction, err := s.transactionRepository.GetByID(context.Background(), transactionID)
	if err != nil {
		return fmt.Errorf("failed to retrieve transaction for sending email: %v", err)
	}

	// Get user associated with the transaction
	user, err := s.transactionRepository.GetUserByID(context.Background(), transaction.IDUser)
	if err != nil {
		return fmt.Errorf("failed to retrieve user for sending email: %v", err)
	}

	// Email template setup
	templatePath := "./templates/email/transaction-success.html"
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %v", err)
	}

	replacerEmail := map[string]interface{}{
		"TransactionID": transaction.IDTransaction,
		"TotalPrice":    transaction.TotalPrice,
		"Quantity":      transaction.QuantityTicket,
		"Status":        transaction.Status,
		"DateOrder":     transaction.DateOrder.Format("2006-01-02 15:04:05"),
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, replacerEmail); err != nil {
		return fmt.Errorf("failed to execute email template: %v", err)
	}

	// Send email to user
	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.SMTPConfig.Email)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Transaction Successful")
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(
		s.cfg.SMTPConfig.Host,
		s.cfg.SMTPConfig.Port,
		s.cfg.SMTPConfig.Email,
		s.cfg.SMTPConfig.Password,
	)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
