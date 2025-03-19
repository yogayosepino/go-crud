package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

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

//update
func UpdateUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == ""{
		http.Error(w, "Id tidak boleh kosong", http.StatusBadRequest)
		return
	}

	var users model.Users
	err := json.NewDecoder(r.Body).Decode(&users)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	switch r.Method{
	case http.MethodPut:
		if users.Username == "" || users.Password == ""{
			http.Error(w, "Semua field harus diisi", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(users.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Gagal meng-hash password", http.StatusInternalServerError)
			return
		}

		_,err = db.Exec("UPDATE users SET username = ?, password = ? WHERE id = ?", users.Username, hashedPassword, id)
		if err != nil {
			fmt.Println("Database error:", err)
			http.Error(w, "Gagal membuat akun, kemungkinan username sudah terdaftar", http.StatusInternalServerError)
			return
		}
		
	case http.MethodPatch:

		var updates []string
		var values []interface{}

		if users.Username != ""{
			updates = append(updates, "username=?")
			values = append(values, users.Username)
		}
		
		if users.Password != ""{
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(users.Password), bcrypt.DefaultCost)
   			if err != nil {
        		http.Error(w, "Gagal meng-hash password", http.StatusInternalServerError)
        		return
    		}
			updates = append(updates, "password=?")
			values = append(values, string(hashedPassword))
		}

		if len(updates) == 0 {
			http.Error(w, "Tidak ada data yang diupdate", http.StatusBadRequest)
            return
		}

		query := fmt.Sprintf("UPDATE users SET %s WHERE id=?", strings.Join(updates, ", "))
        values = append(values, id)
        _, err = db.Exec(query, values...)
    }

    if err != nil {
        http.Error(w, "Gagal mengupdate data", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Data berhasil diperbarui"})
	
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