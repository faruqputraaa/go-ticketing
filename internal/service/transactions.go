package service

import (
	"context"
	_ "errors"
	"fmt"
	"github.com/faruqputraaa/go-ticket/config"
	"github.com/faruqputraaa/go-ticket/internal/entity"
	"github.com/faruqputraaa/go-ticket/internal/http/dto"
	"github.com/faruqputraaa/go-ticket/internal/repository"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"time"
)

type TransactionService interface {
	GetAll(ctx context.Context) ([]entity.Transaction, error)
	GetByID(ctx context.Context, id string) (*entity.Transaction, error)
	GetByIDUser(ctx context.Context, IDUser int) ([]entity.Transaction, error)
	Create(ctx context.Context, req dto.CreateTransactionRequest, claims *entity.JWTCustomClaims) (*entity.Transaction, *snap.Response, error)
	Update(ctx context.Context, transaction dto.UpdateTransactionRequest) error
	LogTransaction(ctx context.Context, transactionID string, status string, message string) error
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

	ticket, err := s.transactionRepository.GetTicketByID(ctx, req.IDTicket)
	if err != nil {
		return nil, nil, fmt.Errorf("Ticket not found")
	}

	// Cek jika harga tiket 0
	if ticket.Price == 0 {
		// Jika harga tiket 0, langsung buat transaksi tanpa Midtrans
		amount := 0.0 // Total harga = 0
		transactionID := fmt.Sprintf("TRX-%d", time.Now().Unix())

		newTransaction := &entity.Transaction{
			IDTransaction:  transactionID,
			IDUser:         claims.IDUser,
			IDTicket:       req.IDTicket,
			QuantityTicket: req.QuantityTicket,
			TotalPrice:     amount,
			Status:         "SUCCESS", // Status langsung sukses karena harga 0
			DateOrder:      time.Now(),
		}

		// Simpan transaksi ke database
		if err := s.transactionRepository.Create(ctx, newTransaction); err != nil {
			return nil, nil, fmt.Errorf("Failed to save transaction")
		}

		// Kembalikan transaksi tanpa response Midtrans
		return newTransaction, nil, nil
	}

	// Jika harga tiket lebih dari 0, lanjutkan dengan proses transaksi melalui Midtrans
	amount := float64(ticket.Price) * float64(req.QuantityTicket)
	transactionID := fmt.Sprintf("TRX-%d", time.Now().Unix())

	newTransaction := &entity.Transaction{
		IDTransaction:  transactionID,
		IDUser:         claims.IDUser,
		IDTicket:       req.IDTicket,
		QuantityTicket: req.QuantityTicket,
		TotalPrice:     amount,
		Status:         "PENDING",
		DateOrder:      time.Now(),
	}

	if err := s.transactionRepository.Create(ctx, newTransaction); err != nil {
		return nil, nil, fmt.Errorf("Failed to save transaction")
	}

	sn := snap.Client{}
	sn.New(s.cfg.MidtransConfig.Serverkey, midtrans.Sandbox)

	reqMidtrans := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  transactionID,
			GrossAmt: int64(amount),
		},
	}

	fmt.Printf("Request to Midtrans: %+v\n", reqMidtrans)

	snapResp, err := sn.CreateTransaction(reqMidtrans)

	fmt.Printf("Response from Midtrans: %+v\n", snapResp)

	return newTransaction, snapResp, nil
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

// Update the Update method
func (s *transactionService) Update(ctx context.Context, req dto.UpdateTransactionRequest) error {
	transaction, err := s.transactionRepository.GetByID(ctx, req.IDTransaction)
	if err != nil {
		return fmt.Errorf("transaction not found: %v", err)
	}

	transaction.Status = req.Status

	return s.transactionRepository.Update(ctx, transaction)

}
