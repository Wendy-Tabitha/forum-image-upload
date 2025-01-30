package handlers

import (
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")

		// Check if the passwords match
		if password != confirmPassword {
			http.Error(w, "Passwords do not match", http.StatusBadRequest)
			return
		}

		// Check if email already exists
		var existingEmail string
		err := db.QueryRow("SELECT email FROM users WHERE email = ?", email).Scan(&existingEmail)
		if err == nil {
			http.Error(w, "Email already taken", http.StatusBadRequest)
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		// Insert new user into the database
		_, err = db.Exec("INSERT INTO users (email, username, password) VALUES (?, ?, ?)", email, username, hashedPassword)
		if err != nil {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/post", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/register.html"))
	tmpl.Execute(w, nil)
}
