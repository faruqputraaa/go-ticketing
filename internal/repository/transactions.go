package repository

import (
	"context"

	"github.com/faruqputraaa/go-ticket/internal/entity"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	GetAll(ctx context.Context) ([]entity.Transaction, error)
	GetByID(ctx context.Context, id string) (*entity.Transaction, error)
	GetByIdUser(ctx context.Context, IDUser int) ([]entity.Transaction, error)
	Create(ctx context.Context, transaction *entity.Transaction) error
	Update(ctx context.Context, transaction *entity.Transaction) error
	GetTicketByID(ctx context.Context, id int64) (*entity.Ticket, error)
	GetUserByID(ctx context.Context, id int) (*entity.User, error)
	CreateLogTransaction(ctx context.Context, log *entity.TransactionLog) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) GetAll(ctx context.Context) ([]entity.Transaction, error) {
	result := make([]entity.Transaction, 0)

	if err := r.db.WithContext(ctx).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *transactionRepository) GetByID(ctx context.Context, id string) (*entity.Transaction, error) {
	result := new(entity.Transaction)

	if err := r.db.WithContext(ctx).Where("id_transaction = ?", id).First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *transactionRepository) Create(ctx context.Context, transaction *entity.Transaction) error {
	return r.db.WithContext(ctx).Create(transaction).Error
}

func (r *transactionRepository) GetTicketByID(ctx context.Context, id int64) (*entity.Ticket, error) {
	result := new(entity.Ticket)

	if err := r.db.WithContext(ctx).Where("id_ticket = ?", id).First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *transactionRepository) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	result := new(entity.User)

	if err := r.db.WithContext(ctx).Where("id_user = ?", id).First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *transactionRepository) Update(ctx context.Context, transaction *entity.Transaction) error {
	return r.db.WithContext(ctx).
		Where("id_transaction = ?", transaction.IDTransaction).
		Updates(transaction).Error
}

func (r *transactionRepository) GetByIdUser(ctx context.Context, IDUser int) ([]entity.Transaction, error) {
	var transaction []entity.Transaction
	err := r.db.WithContext(ctx).Where("id_user = ?", IDUser).Find(&transaction).Error
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (r *transactionRepository) CreateLogTransaction(ctx context.Context, log *entity.TransactionLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}
