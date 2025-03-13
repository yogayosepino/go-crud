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


//POST
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

//Update PUT
func UpdateEmployee(db *sql.DB, w http.ResponseWriter, r *http.Request){
	if r.Method != "PUT"{
		http.Error(w, `{"Method Not Allowed}`, http.StatusMethodNotAllowed)
		return
	}

	//ambil id
	id := r.URL.Query().Get("id")
	if id == ""{
		http.Error(w, `{"error" : "Id tidak boleh kosong}`, http.StatusBadRequest)
		return
	}

	//decode json req body
	var updatedData model.Employee
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, `{"error" : "format json tidak valid"}`, http.StatusBadRequest)
		return
	}

	//validasi field tidak boleh kosong
	if updatedData.Name == "" || updatedData.NPWP == "" || updatedData.Address == "" {
		http.Error(w, `{"error" : "field tidak boleh kosong"}`, http.StatusBadRequest)
		return
	}

	//id checking apakah ada di db
	var existing model.Employee
	err := db.QueryRow("SELECT id, name, npwp, address FROM employee WHERE id = ?", id).Scan(&existing.Id,&existing.Name, &existing.NPWP, &existing.Address)

	if err == sql.ErrNoRows{
		http.Error(w, `{"error" : "Data tidak ditemukan"}`, http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, `{"error" : "Gagal mengambil data"}`, http.StatusInternalServerError)
		return
	}

	//update data
	_, err = db.Exec("UPDATE employee SET name = ?, npwp = ?, address = ? WHERE id = ?", updatedData.Name, updatedData.NPWP, updatedData.Address, id)
	if err != nil{
		http.Error(w, `{"error" : "Gagal mengupdate data "}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedData)
}


