package entity

import "time"

type Event struct {
	IDEvent     int64     `json:"id"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	Time        time.Time `json:"time"`
	Description string    `json:"description"`
}

func (Event) TableName() string {
	return "events"
}
