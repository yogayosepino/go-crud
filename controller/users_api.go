package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	//"io"
	"net/http"

	"github.com/yogayosepino/go-crud/model"
	"golang.org/x/crypto/bcrypt"
)

//get
func GetUsers(db *sql.DB) ([]model.UserResponse, error) {
	var users []model.UserResponse

	// Coba jalankan query manual
	rows, err := db.Query("SELECT id, username FROM users")
	if err != nil {
		fmt.Println("Query error:", err)
		return nil, err
	}
	defer rows.Close()

	// Debug: Cek apakah ada baris hasil dari query
	count := 0
	for rows.Next() {
		var user model.UserResponse
		if err := rows.Scan(&user.Id, &user.Username); err != nil {
			fmt.Println("Scan error:", err)
			return nil, err
		}
		users = append(users, user)
		count++
	}

	if count == 0 {
		fmt.Println("Query sukses tapi tidak ada data!")
		return nil, fmt.Errorf("tidak ada data")
	}

	fmt.Println("Data berhasil diambil:", users)
	return users, nil
}


//post
func CreateUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var users model.Users

	// Decode JSON langsung
	err := json.NewDecoder(r.Body).Decode(&users)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	// Debugging: cek apakah username ada isinya
	fmt.Println("Username:", users.Username)
	fmt.Println("Password:", users.Password)

	// Validasi username tidak boleh kosong
	if users.Username == "" {
		http.Error(w, "Username tidak boleh kosong", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(users.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Gagal membuat akun", http.StatusInternalServerError)
		return
	}

	// Insert ke database
	result, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", users.Username, string(hashedPassword))
	if err != nil {
		fmt.Println("Database error:", err)
		http.Error(w, "Gagal membuat akun, kemungkinan username sudah terdaftar", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Println("Rows affected:", rowsAffected)

	// Response sukses
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Users created successfully"})
}

//delete
func DeleteUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE"{
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error": "Id tidak boleh kosong"}`, http.StatusBadRequest)
		return
	}

	_,  err := db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Failed to delete data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Users deleted succesfully"})

}