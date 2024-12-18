package entity

import "time"

type Transaction struct {
	IDTransaction  int64     `json:"id_transaction"`
	IDUser         int       `json:"id_user"`
	QuantityTicket int       `json:"quantity_ticket"`
	IDTicket       int       `json:"id_ticket"`
	TotalPrice     float64   `json:"total_price"`
	DateOrder      time.Time `json:"date_order"`
}

func (Transaction) TableName() string {
	return "transaction"
}
