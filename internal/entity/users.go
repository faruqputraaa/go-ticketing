package entity

type User struct {
	IDUser   int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Role     string `json:"role"`
}

func (User) TableName() string {
	return "users"
}
