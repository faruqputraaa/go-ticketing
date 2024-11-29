package entity

type Ticket struct {
	IDTicket int64  `json:"id_ticket"`
	Price    int64  `json:"price"`
	Category string `json:"category"`
	IDEvent  int64  `json:"id_event"`
}

func (Ticket) TableName() string {
	return "ticket"
}
