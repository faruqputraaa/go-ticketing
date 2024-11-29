package repository

import (
	"context"

	"github.com/faruqputraaa/go-ticket/internal/entity"
	"gorm.io/gorm"
)

type TicketRepository interface {
	GetAll(ctx context.Context) ([]entity.Ticket, error)
	GetByID(ctx context.Context, id int64) (*entity.Ticket, error)
	Create(ctx context.Context, ticket *entity.Ticket) error
	Update(ctx context.Context, ticket *entity.Ticket) error
	Delete(ctx context.Context, ticket *entity.Ticket) error
}

type ticketRepository struct {
	db *gorm.DB
}

// Create implements TicketRepository.
func (r *ticketRepository) Create(ctx context.Context, ticket *entity.Ticket) error {
	return r.db.WithContext(ctx).Create(ticket).Error
}

// Delete implements TicketRepository.
func (r *ticketRepository) Delete(ctx context.Context, ticket *entity.Ticket) error {
	return r.db.WithContext(ctx).Delete(ticket).Error
}

// GetAll implements TicketRepository.
func (r *ticketRepository) GetAll(ctx context.Context) ([]entity.Ticket, error) {
	result := make([]entity.Ticket, 0)

	if err := r.db.WithContext(ctx).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

// GetByID implements TicketRepository.
func (r *ticketRepository) GetByID(ctx context.Context, id int64) (*entity.Ticket, error) {
	result := new(entity.Ticket)
	if err := r.db.WithContext(ctx).First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

// Update implements TicketRepository.
func (r *ticketRepository) Update(ctx context.Context, ticket *entity.Ticket) error {
	return r.db.WithContext(ctx).Updates(ticket).Error
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{db}
}
