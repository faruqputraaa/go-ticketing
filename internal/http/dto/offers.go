package dto

type GetOfferByIDRequest struct {
	IDOffer int64 `param:"id_offer" validate:"required"`
}

type GetOfferByIDUserRequest struct {
	IDUser int64 `json:"id_user" validate:"required"`
}

type CreateOfferRequest struct {
	IDUser      int64  `json:"id_user" validate:"required"`
	Email       string `json:"email" validate:"required"`
	NameEvent   string `json:"name_event" validate:"required"`
	Description string `json:"description" validate:"required"`
	Status      string `json:"status" validate:"required"`
}

type UpdateOfferRequest struct {
	IDOffer     int64  `param:"id_offer" validate:"required"`
	IDUser      int64  `json:"id_user" validate:"required"`
	Email       string `json:"email" validate:"required"`
	NameEvent   string `json:"name_event" validate:"required"`
	Description string `json:"description" validate:"required"`
	Status      string `json:"status" validate:"required"`
}
