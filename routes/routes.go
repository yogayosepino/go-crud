package routes

import (
	"database/sql"
	"encoding/json"
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

	server.HandleFunc("/api/employees", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method{
		case http.MethodGet :
				employees, err := controller.GetEmployees(db)
				if err != nil {
					http.Error(w, "Gagal mengambil data", http.StatusInternalServerError)
					return
				}
	
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(employees)
		case http.MethodPost :
			controller.CreateEmployee(db,w,r)

		default:
			http.Error(w, "Method Now Allowed", http.StatusMethodNotAllowed)
		}
		
	})

	server.HandleFunc("/api/employees/update", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method{
		case http.MethodPut :
			controller.UpdateEmployee(db,w,r)
		
		case http.MethodPatch:
			controller.UpdateEmployeePatch(db,w,r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
		
	})

	server.HandleFunc("/api/employees/delete", func(w http.ResponseWriter, r *http.Request){
		if r.Method == http.MethodDelete{
			controller.DeleteEmployee(db,w,r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})


	server.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request){
		switch r.Method{
		case http.MethodGet :
			users, err := controller.GetUsers(db)
				if err != nil {
					http.Error(w, "Gagal mengambil data", http.StatusInternalServerError)
					return
				}
	
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(users)
				
		default :
			http.Error(w, "Method Now Allowed", http.StatusMethodNotAllowed)

		}

	})
}

