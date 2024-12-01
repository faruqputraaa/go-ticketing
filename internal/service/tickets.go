package service

import (
	"context"

	"github.com/faruqputraaa/go-ticket/internal/entity"
	"github.com/faruqputraaa/go-ticket/internal/repository"
)

type TicketService interface {
	GetAll(ctx context.Context) ([]entity.Ticket, error)
	GetByID(ctx context.Context, id int64) (*entity.Ticket, error)
	Create(ctx context.Context, ticket *entity.Ticket) error
	Update(ctx context.Context, ticket *entity.Ticket) error
	Delete(ctx context.Context, ticket *entity.Ticket) error
}

type ticketService struct {
	ticketRepository repository.TicketRepository
}

func NewTicketService(ticketRepository repository.TicketRepository) TicketService {
	return &ticketService{ticketRepository}
}

// Create implements TicketService.
func (s *ticketService) Create(ctx context.Context, ticket *entity.Ticket) error {
	return s.ticketRepository.Create(ctx, ticket)
}

// Delete implements TicketService.
func (s *ticketService) Delete(ctx context.Context, ticket *entity.Ticket) error {
	return s.ticketRepository.Delete(ctx, ticket)
}

// GetAll implements TicketService.
func (s *ticketService) GetAll(ctx context.Context) ([]entity.Ticket, error) {
	return s.ticketRepository.GetAll(ctx)
}

// GetByID implements TicketService.
func (s *ticketService) GetByID(ctx context.Context, id int64) (*entity.Ticket, error) {
	return s.ticketRepository.GetByID(ctx, id)
}

// Update implements TicketService.
func (s *ticketService) Update(ctx context.Context, ticket *entity.Ticket) error {
	return s.ticketRepository.Update(ctx, ticket)

}
