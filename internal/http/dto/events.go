package dto

import "time"

type GetEventByIDRequest struct {
	IDEvent int64 `param:"id_event" validate:"required"`
}

type CreateEventRequest struct {
	Name        string    `json:"name" validate:"required"`
	Location    string    `json:"location" validate:"required"`
	Time        time.Time `json:"time" validate:"required"`
	Description string    `json:"description" validate:"required"`
}

type UpdateEventRequest struct {
	IDEvent     int64     `param:"id_event" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	Location    string    `json:"location" validate:"required"`
	Time        time.Time `json:"time" validate:"required"`
	Description string    `json:"description" validate:"required"`
}
