package model

type Users struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct{
	Id 		 string `json:"id"`
	Username string `json:"username"`
}