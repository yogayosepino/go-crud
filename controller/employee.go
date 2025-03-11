package controller

import (
	"database/sql"
	"fmt"

	"github.com/yogayosepino/go-crud/model"
)

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

