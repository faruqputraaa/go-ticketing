package entity

import "time"

type Event struct {
	IDEvent     int64     `json:"id_event" gorm:"primary_key;auto_increment"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	Time        time.Time `json:"time"`
	Description string    `json:"description"`
}

func (Event) TableName() string {
	return "events"
}
