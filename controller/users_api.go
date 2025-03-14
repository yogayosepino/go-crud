package controller

import (
	"database/sql"
	"fmt"

	"github.com/yogayosepino/go-crud/model"
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