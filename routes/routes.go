package routes

import (
	"database/sql"
	"net/http"

	"github.com/yogayosepino/go-crud/controller"
	"github.com/yogayosepino/go-crud/middleware"
)

func MapRoutes(server *http.ServeMux, db *sql.DB) {
	server.HandleFunc("/", controller.NewHelloWorldController())

	server.HandleFunc("/employee", middleware.AuthMiddleware(controller.NewIndexEmployeeController(db)))
	server.HandleFunc("/employee/create", middleware.AuthMiddleware(controller.NewCreateEmployeeController(db)))
	server.HandleFunc("/employee/update", middleware.AuthMiddleware(controller.NewUpdateEmployeeController(db)))
	server.HandleFunc("/employee/delete", middleware.AuthMiddleware(controller.NewDeleteEmployeeController(db)))

	server.HandleFunc("/login", controller.NewLoginController(db))
	server.HandleFunc("/register", controller.NewSignupController(db))

}
