package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("secret-key"))

func init() {
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400, // 1 hari
		HttpOnly: true,
	}
}

func AuthMiddleware(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "session") // Pastikan nama session sama
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		auth, ok := session.Values["authenticated"].(bool)
		if !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		handlerFunc(w, r)
	}
}
