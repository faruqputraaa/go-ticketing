package entity

import "time"

type Transaction struct {
	IDTransaction  int       `json:"id_transaction"`
	IDUser         int       `json:"id_user"`
	QuantityTicket int       `json:"quantity_ticket"`
	IDEvent        int       `json:"id_event"`
	IDTicket       int       `json:"id_ticket"`
	TotalPrice     float64   `json:"total_price"`
	DateOrder      time.Time `json:"date_order"`
	CreatedAt      time.Time `json:"created_at"`
}

func (Transaction) TableName() string {
	return "transaction"
}
