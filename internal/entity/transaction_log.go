package entity

import "time"

type TransactionLog struct {
	IDTransactionLog int       `json:"id_transaction_log" gorm:"primaryKey"`
	IDTransaction    string    `gorm:"column:id_transaction"`
	Status           string    `json:"status"`
	Message          string    `json:"message"`
	CreatedAt        time.Time `json:"created_at"`
}

func (TransactionLog) TableName() string {
	return "transaction_log"
}
