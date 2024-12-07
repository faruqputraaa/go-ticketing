package entity

type User struct {
	IDUser   int    `json:"id_user"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Role     string `json:"role"`
}

func (User) TableName() string {
	return "users"
}
