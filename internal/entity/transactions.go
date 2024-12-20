package entity

import "time"

type Transaction struct {
	IDTransaction  string    `json:"id_transaction"`
	IDUser         int       `json:"id_user"`
	QuantityTicket int       `json:"quantity_ticket"`
	IDTicket       int64     `json:"id_ticket"`
	TotalPrice     float64   `json:"total_price"`
	Status         string    `json:"status"`
	DateOrder      time.Time `json:"date_order"`
}

func (Transaction) TableName() string {
	return "transactions"
}
