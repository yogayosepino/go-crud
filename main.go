package main

import (
	//"fmt"
	"net/http"

	//"github.com/yogayosepino/go-crud/controller"
	"github.com/yogayosepino/go-crud/database"
	"github.com/yogayosepino/go-crud/routes"
)

func main() {
	db := database.InitDatabase()

	// fmt.Println("Hello")

	server := http.NewServeMux()

	routes.MapRoutes(server, db)

	http.ListenAndServe(":8080", server)

	//testing
}
