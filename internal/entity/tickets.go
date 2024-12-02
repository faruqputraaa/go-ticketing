package entity

type Ticket struct {
	IDTicket int64  `json:"id_ticket"`
	IDEvent  int64  `json:"id_event"`
	Price    int64  `json:"price"`
	Category string `json:"category"` //jenis tiket

}

func (Ticket) TableName() string {
	return "ticket"
}
