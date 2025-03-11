package controller

import (
	"database/sql"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gorilla/sessions"
	"github.com/yogayosepino/go-crud/model"
	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte("secret-key"))

// type Users struct {
// 	Id       int
// 	Username string
// 	Password string
// }

func NewLoginController(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			password := r.FormValue("password")

			var users model.Users
			err := db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username).
				Scan(&users.Id, &users.Username, &users.Password)
			if err != nil {
				http.Error(w, "Username atau password salah", http.StatusUnauthorized)
				return
			}

			err = bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(password))
			if err != nil {
				http.Error(w, "Username atau password salah", http.StatusUnauthorized)
				return
			}

			session, _ := store.Get(r, "session")
			session.Values["authenticated"] = true
			session.Values["username"] = users.Username
			session.Save(r, w)

			http.Redirect(w, r, "/employee", http.StatusSeeOther)
			return
		}

		fp := filepath.Join("views", "login.html")
		tmpl, err := template.ParseFiles(fp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, nil)
	}
}

func NewSignupController(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			password := r.FormValue("password")

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "Gagal membuat akun", http.StatusInternalServerError)
				return
			}

			_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, hashedPassword)
			if err != nil {
				http.Error(w, "Gagal membuat akun", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		fp := filepath.Join("views", "register.html")
		tmpl, err := template.ParseFiles(fp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, nil)
	}
}
