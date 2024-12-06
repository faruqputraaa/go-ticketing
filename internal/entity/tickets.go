package entity

type Ticket struct {
	IDTicket int64  `json:"id_ticket" gorm:"primary_key;auto_increment"`
	IDEvent  int64  `json:"id_event"`
	Price    int64  `json:"price"`
	Category string `json:"category"` //jenis tiket

}

func (Ticket) TableName() string {
	return "tickets"
}
