package service

import (
	"context"

	"github.com/faruqputraaa/go-ticket/internal/entity"
	"github.com/faruqputraaa/go-ticket/internal/http/dto"
	"github.com/faruqputraaa/go-ticket/internal/repository"
)

type EventService interface {
	Create(ctx context.Context, req dto.CreateEventRequest) error
	GetAll(ctx context.Context, req dto.GetAllEventsRequest) ([]entity.Event, error)
	GetByID(ctx context.Context, id int64) (*entity.Event, error)
	Update(ctx context.Context, req dto.UpdateEventRequest) error
	Delete(ctx context.Context, event *entity.Event) error
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

func (s *eventService) GetAll(ctx context.Context, req dto.GetAllEventsRequest) ([]entity.Event, error) {
	return s.eventRepository.GetAll(ctx, req)
}

func (s *eventService) GetByID(ctx context.Context, id int64) (*entity.Event, error) {
	return s.eventRepository.GetByID(ctx, id)
}

func (s *eventService) Update(ctx context.Context, req dto.UpdateEventRequest) error {
	event, err := s.eventRepository.GetByID(ctx, req.IDEvent)
	if err != nil {
		return err
	}
	if req.IDEvent != 0 {
		event.IDEvent = req.IDEvent
	}
	if req.Name != "" {
		event.Name = req.Name
	}
	if req.Location != "" {
		event.Location = req.Location
	}

	if req.Description != "" {
		event.Description = req.Description
	}
	return s.eventRepository.Update(ctx, event)
}

func (s *eventService) Delete(ctx context.Context, event *entity.Event) error {
	return s.eventRepository.Delete(ctx, event)
}
