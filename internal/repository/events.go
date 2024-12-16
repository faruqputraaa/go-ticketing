package repository

import (
	"context"

	"github.com/faruqputraaa/go-ticket/internal/entity"
	"gorm.io/gorm"
)

type EventRepository interface {
	GetAll(ctx context.Context) ([]entity.Event, error)
	GetByID(ctx context.Context, id int64) (*entity.Event, error)
	Create(ctx context.Context, event *entity.Event) error
	Update(ctx context.Context, event *entity.Event) error
	Delete(ctx context.Context, event *entity.Event) error
	SearchByName(ctx context.Context, name string) ([]entity.Event, error)
	SearchByLocation(ctx context.Context, location string) ([]entity.Event, error)
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db}
}

// GetAll
func (r *eventRepository) GetAll(ctx context.Context) ([]entity.Event, error) {
	result := make([]entity.Event, 0)

	if err := r.db.WithContext(ctx).Find(&result).Error; err != nil {
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

func (r *eventRepository) SearchByName(ctx context.Context, name string) ([]entity.Event, error) {
	var events []entity.Event
	if err := r.db.WithContext(ctx).Where("name LIKE ?", "%"+name+"%").Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (r *eventRepository) SearchByLocation(ctx context.Context, location string) ([]entity.Event, error) {
	var events []entity.Event
	if err := r.db.WithContext(ctx).Where("location LIKE ?", "%"+location+"%").Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}
