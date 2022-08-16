package model

// User data
type User struct {
	ID       int    `json:"id"`
	UserName string `json:"username" form:"username" binding:"required"`
	Password string `json:"password,omitempty" form:"password" binding:"required"`
}
