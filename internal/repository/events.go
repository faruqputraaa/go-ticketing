package repository

import (
	"context"
	"strings"

	"github.com/faruqputraaa/go-ticket/internal/entity"
	"github.com/faruqputraaa/go-ticket/internal/http/dto"
	"gorm.io/gorm"
)

type EventRepository interface {
	GetAll(ctx context.Context, req dto.GetAllEventsRequest) ([]entity.Event, error)
	GetByID(ctx context.Context, id int64) (*entity.Event, error)
	Create(ctx context.Context, event *entity.Event) error
	Update(ctx context.Context, event *entity.Event) error
	Delete(ctx context.Context, event *entity.Event) error
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db}
}

// GetAll
func (r *eventRepository) GetAll(ctx context.Context, req dto.GetAllEventsRequest) ([]entity.Event, error) {
	result := make([]entity.Event, 0)

	query := r.db.WithContext(ctx)
	if req.Search != "" {
		search := strings.ToLower(req.Search)
		query = query.Where("LOWER(name) LIKE ?", "%"+search+"%").
			Or("LOWER(location) LIKE ?", "%"+search+"%")
	}

	if req.Sort != "" && req.Order != "" {
		query = query.Order(req.Sort + " " + req.Order)
	}

	if req.Page > 0 && req.Limit > 0 {
		query = query.Offset((req.Page - 1) * req.Limit).Limit(req.Limit)
	}

	if err := query.Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

// GetByID
func (r *eventRepository) GetByID(ctx context.Context, id int64) (*entity.Event, error) {
	result := new(entity.Event)
	if err := r.db.WithContext(ctx).Where("id_event = ?", id).First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

// Create
func (r *eventRepository) Create(ctx context.Context, event *entity.Event) error {
	return r.db.WithContext(ctx).Create(event).Error
}

// Update
func (r *eventRepository) Update(ctx context.Context, event *entity.Event) error {
	return r.db.WithContext(ctx).Updates(event).Error
}

// Delete
func (r *eventRepository) Delete(ctx context.Context, event *entity.Event) error {
	return r.db.WithContext(ctx).Where("id_event = ?", event.IDEvent).Delete(&event).Error
}
