package model

type Users struct {
	Id       string `json:"id"`
	Username string `json:"name"`
	Password string `json:"password"`
}

type UserResponse struct{
	Id 		 string `json:"id"`
	Username string `json:"name"`
}