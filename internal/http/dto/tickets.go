package dto

type GetTicketByIDRequest struct {
	IDTicket int64 `param:"id" validate:"required"`
}

type GetTicketByIDEventRequest struct {
	IDEvent int64 `json:"id_event,string" validate:"required"`
}

type CreateTicketRequest struct {
	IDEvent  int64  `json:"id_event,string" validate:"required"`
	Price    int64  `json:"price,string" validate:"required"`
	Category string `json:"category" validate:"required"`
}

type UpdateTicketRequest struct {
	IDTicket int64  `param:"id" validate:"required"`
	IDEvent  int64  `json:"id_event,string" validate:"required"`
	Price    int64  `json:"price" validate:"required"`
	Category string `json:"category" validate:"required"`
}
