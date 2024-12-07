package entity

type Offer struct {
	IDOffer     int64  `json:"id_offer" gorm:"autoIncrement;"`
	IDUser      int64  `json:"id_user"`
	Email       string `json:"email"`
	NameEvent   string `json:"name_event"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
