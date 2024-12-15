package entity

type User struct {
	IDUser             int    `json:"id_user" gorm:"autoIncrement;"`
	Username           string `json:"username"`
	Email              string `json:"email"`
	Password           string `json:"-"`
	Role               string `json:"role"`
	ResetPasswordToken string `json:"reset_password_token"`
	VerifyEmailToken   string `json:"verify_email_token"`
	IsVerified         int    `json:"is_verified"`
}

func (User) TableName() string {
	return "users"
}
