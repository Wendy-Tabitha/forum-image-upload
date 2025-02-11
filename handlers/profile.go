package handlers

import (
    "encoding/json"
    "log"
    "net/http"
)

// ProfileHandler handles the profile page
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
    userID := getUserIDFromSession(w, r)
    if userID == "" {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    // Fetch user details from the database
    var user User
    err := db.QueryRow("SELECT id, email, username FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Email, &user.Username)
    if err != nil {
        log.Printf("Error fetching user details: %v", err)
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }

    // Return user details as JSON
    w.Header().Set("Content-Type", "application/json")
    err = json.NewEncoder(w).Encode(user)
    if err != nil {
        log.Printf("Error encoding user details: %v", err)
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }
}