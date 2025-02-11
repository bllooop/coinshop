package domain

type User struct {
	Id       int    `json:"-" db:"id"`
	UserName string `json:"username"`
	Password string `json:"password"`
	Coins    *int   `json:"coins"`
}

type SignInInput struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
