package model

// LoginData contains the data used for login
type LoginData struct {
	UserName string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Room     string `json:"room" form:"room" binding:"required"`
}
