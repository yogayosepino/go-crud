package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yogayosepino/go-crud/model"
)

//GET
func GetEmployees(db *sql.DB) ([]model.Employee, error) {
	var employees []model.Employee

	// Coba jalankan query manual
	rows, err := db.Query("SELECT id, name, npwp, address FROM employee")
	if err != nil {
		fmt.Println("Query error:", err)
		return nil, err
	}
	defer rows.Close()

	// Debug: Cek apakah ada baris hasil dari query
	count := 0
	for rows.Next() {
		var emp model.Employee
		if err := rows.Scan(&emp.Id, &emp.Name, &emp.NPWP, &emp.Address); err != nil {
			fmt.Println("Scan error:", err)
			return nil, err
		}
		employees = append(employees, emp)
		count++
	}

	if count == 0 {
		fmt.Println("Query sukses tapi tidak ada data!")
		return nil, fmt.Errorf("tidak ada data")
	}

	fmt.Println("Data berhasil diambil:", employees)
	return employees, nil
}

func CreateEmployee(db *sql.DB, w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var employee model.Employee

	err := json.NewDecoder(r.Body).Decode(&employee)

	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO employee (name, npwp, address) VALUES (?,?,?)"
	_, err = db.Exec(query, employee.Name, employee.NPWP, employee.Address)
	if err != nil {
		http.Error(w, "Failed to insert data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Employee created succesfully"})

}

