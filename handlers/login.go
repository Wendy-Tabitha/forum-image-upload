package handlers

import (
	"html/template"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			http.Error(w, "Error parsing file", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	} else if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		var user User
		err := db.QueryRow("SELECT id, password FROM users WHERE email = ?", email).Scan(&user.ID, &user.Password)
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// Set session cookie
		http.SetCookie(w, &http.Cookie{
			Name:    "session_id",
			Value:   "some_unique_session_id",       // You should generate a unique session ID
			Expires: time.Now().Add(24 * time.Hour), // Set expiration
		})

		http.Redirect(w, r, "/post", http.StatusSeeOther)
		return
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}
