package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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

//update patch
// Update PATCH
func UpdateEmployeePatch(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != "PATCH" {
		http.Error(w, `{"error": "Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error": "Id tidak boleh kosong"}`, http.StatusBadRequest)
		return
	}

	fmt.Println("ID yang diterima:", id) // DEBUG 1: Cek ID

	var updates map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Format JSON salah", http.StatusBadRequest)
		return
	}

	// Jika tidak ada field yang dikirim, tolak request
	if len(updates) == 0 {
		http.Error(w, "tidak ada data yang diupdate", http.StatusBadRequest)
		return
	}

	// Build query update secara dinamis
	var setClauses []string
	var values []interface{}
	for key, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = ?", key))
		values = append(values, value)
	}

	query := fmt.Sprintf("UPDATE employee SET %s WHERE id = ?", strings.Join(setClauses, ", "))
	values = append(values, id)

	fmt.Println("QUERY:", query) // DEBUG 2: Cek Query yang Dibentuk
	fmt.Println("VALUES:", values) // DEBUG 3: Cek Nilai yang akan Dimasukkan

	// Eksekusi query
	res, err := db.Exec(query, values...)
	if err != nil {
		fmt.Println("ERROR:", err) // DEBUG 4: Print Error dari Database
		http.Error(w, "Gagal mengupdate data", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	fmt.Println("ROWS AFFECTED:", rowsAffected) // DEBUG 5: Cek Jumlah Baris yang Terupdate

	if rowsAffected == 0 {
		http.Error(w, "Data tidak ditemukan atau tidak ada perubahan", http.StatusNotFound)
		return
	}

	// Beri respon sukses
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Data berhasil diperbarui"})
}


