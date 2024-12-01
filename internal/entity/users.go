package entity

type User struct {
	IDUser   int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (User) TableName() string {
	return "user"
}
