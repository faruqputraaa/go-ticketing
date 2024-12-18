package service

import (
	"context"
	"errors"
	_ "errors"
	"github.com/faruqputraaa/go-ticket/config"
	"github.com/faruqputraaa/go-ticket/internal/entity"
	"github.com/faruqputraaa/go-ticket/internal/http/dto"
	"github.com/faruqputraaa/go-ticket/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type TransactionService interface {
	GetAll(ctx context.Context) ([]entity.Transaction, error)
	GetByID(ctx context.Context, id int64) (*entity.Transaction, error)
	GetByIDUser(ctx context.Context, IDUser int) ([]entity.Transaction, error)
	Create(ctx echo.Context, req dto.CreateTransactionRequest) error
	Update(ctx context.Context, transaction dto.UpdateTransactionRequest) error
}

type transactionService struct {
	cfg                   *config.Config
	transactionRepository repository.TransactionRepository
}

func NewTransactionService(cfg *config.Config, transactionRepository repository.TransactionRepository) TransactionService {
	return &transactionService{cfg, transactionRepository}
}

func (s *transactionService) Create(ctx echo.Context, req dto.CreateTransactionRequest) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*entity.JWTCustomClaims)

	if req.QuantityTicket == 0 {
		return errors.New("email is required")
	}

	transaction := &entity.Transaction{
		IDUser:         claims.IDUser,
		IDTicket:       req.IDTicket,
		QuantityTicket: req.QuantityTicket,
	}

	var sn = snap.Client{}
	sn.New(s.cfg.MidtransConfig.Serverkey, midtrans.Sandbox)

	return s.transactionRepository.Create(ctx.Request().Context(), transaction)
}

// GetAll implements TicketService.
func (s *transactionService) GetAll(ctx context.Context) ([]entity.Transaction, error) {
	return s.transactionRepository.GetAll(ctx)
}

// GetByID implements TicketService.
func (s *transactionService) GetByID(ctx context.Context, id int64) (*entity.Transaction, error) {
	return s.transactionRepository.GetByID(ctx, id)
}

func (s *transactionService) GetByIDUser(ctx context.Context, IDUser int) ([]entity.Transaction, error) {
	return s.transactionRepository.GetByIdUser(ctx, IDUser)
}

func (s *transactionService) Update(ctx context.Context, req dto.UpdateTransactionRequest) error {
	transaction, err := s.transactionRepository.GetByID(ctx, req.IDTransaction)
	if err != nil {
		return err
	}

	return s.transactionRepository.Update(ctx, transaction)
}
