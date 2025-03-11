package model

type Employee struct {
	Id string `json:"id"`
	Name    string `json:"name"`
	NPWP    string `json:"npwp"`
	Address string `json:"address"`
}