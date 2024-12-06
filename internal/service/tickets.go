package service

import (
	"context"
	_ "errors"
	"fmt"

	"github.com/faruqputraaa/go-ticket/internal/entity"
	"github.com/faruqputraaa/go-ticket/internal/http/dto"
	"github.com/faruqputraaa/go-ticket/internal/repository"
)

type TicketService interface {
	GetAll(ctx context.Context) ([]entity.Ticket, error)
	GetByID(ctx context.Context, id int64) (*entity.Ticket, error)
	GetByIDEvent(ctx context.Context, IDevent int64) ([]entity.Ticket, error)
	Create(ctx context.Context, req dto.CreateTicketRequest) error
	Update(ctx context.Context, req dto.UpdateTicketRequest) error
	Delete(ctx context.Context, ticket *entity.Ticket) error
}

type ticketService struct {
	ticketRepository repository.TicketRepository
}

func NewTicketService(ticketRepository repository.TicketRepository) TicketService {
	return &ticketService{ticketRepository}
}

// Create implements TicketService.
func (s *ticketService) Create(ctx context.Context, req dto.CreateTicketRequest) error {
	ticket := &entity.Ticket{
		IDEvent:  req.IDEvent,
		Price:    req.Price,
		Category: req.Category,
	}

	// Periksa duplikasi berdasarkan IDEvent dan Category
	existingTickets, err := s.ticketRepository.GetByIdEvent(ctx, ticket.IDEvent)
	if err != nil {
		return err
	}
	for _, t := range existingTickets {
		if t.Category == ticket.Category {
			return fmt.Errorf("ticket with category '%s' for event ID %d already exists", ticket.Category, ticket.IDEvent)
		}
	}
	fmt.Printf("Creating ticket: %+v\n", ticket)
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

// GetByIDEvent implements TicketService.
func (s *ticketService) GetByIDEvent(ctx context.Context, IDevent int64) ([]entity.Ticket, error) {
	return s.ticketRepository.GetByIdEvent(ctx, IDevent)
}

// Update implements TicketService.
func (s *ticketService) Update(ctx context.Context, req dto.UpdateTicketRequest) error {
	ticket, err := s.ticketRepository.GetByID(ctx, req.IDTicket)
	if err != nil {
		return err
	}
	if req.IDEvent != 0 {
		ticket.IDEvent = req.IDEvent
	}
	if req.Price != 0 {
		ticket.Price = req.Price
	}
	if req.Category != "" {
		ticket.Category = req.Category
	}
	return s.ticketRepository.Update(ctx, ticket)
}
