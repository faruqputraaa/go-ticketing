package entity

type Ticket struct {
	IDEvent  int64  `json:"id_event"`
	Price    int64  `json:"price"`
	Category string `json:"category"` //jenis tiket

}

func (Ticket) TableName() string {
	return "tickets"
}
