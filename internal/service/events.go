package service

import (
	"context"
	"github.com/faruqputraaa/go-ticket/internal/entity"
	"github.com/faruqputraaa/go-ticket/internal/http/dto"
	"github.com/faruqputraaa/go-ticket/internal/repository"
)

type EventService interface {
	Create(ctx context.Context, req dto.CreateEventRequest) error
	GetAll(ctx context.Context) ([]entity.Event, error)
}

type eventService struct {
	eventRepository repository.EventRepository
}

func NewEventService(eventRepository repository.EventRepository) EventService {
	return &eventService{eventRepository}
}

func (s *eventService) Create(ctx context.Context, req dto.CreateEventRequest) error {
	event := &entity.Event{
		Name:        req.Name,
		Time:        req.Time,
		Location:    req.Location,
		Description: req.Description,
	}

	return s.eventRepository.Create(ctx, event)
}

func (s *eventService) GetAll(ctx context.Context) ([]entity.Event, error) {
	return s.eventRepository.GetAll(ctx)
}
