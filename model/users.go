package model

type Users struct {
	Id       int
	Username string
	Password string
}

type UserResponse struct{
	Id int
	Username string
}