package dto

type GetTransactionByIDRequest struct {
	IDTransaction int64 `param:"id_transaction" validate:"required"`
}

type GetTransactionByIDUserRequest struct {
	IDUser int `json:"id_user" validate:"required"`
}

type CreateTransactionRequest struct {
	QuantityTicket int     `json:"quantity_ticket" validate:"required"`
	IDTicket       int     `json:"id_ticket" validate:"required"`
	TotalPrice     float64 `json:"total_price" validate:"required"`
}

type UpdateTransactionRequest struct {
	IDTransaction  int64  `param:"id_offer" validate:"required"`
	IDUser         int    `json:"id_user" validate:"required"`
	QuantityTicket string `json:"quantity_ticket" validate:"required"`
	IDTicket       string `json:"id_ticket" validate:"required"`
	TotalPrice     string `json:"description" validate:"required"`
	Status         string `json:"status" validate:"required"`
}
