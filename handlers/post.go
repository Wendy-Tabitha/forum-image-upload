package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path == "/post" {
	// }

	if r.Method != http.MethodPost {
		RenderError(w, r, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Check if the user is logged in
	var userID string
	sessionCookie, err := r.Cookie("session_id")
	if err == nil {
		err = db.QueryRow("SELECT user_id FROM sessions WHERE session_id = ?", sessionCookie.Value).Scan(&userID)
		if err == sql.ErrNoRows {
			// Clear the invalid session cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    "",
				Path:     "/",
				Expires:  time.Unix(0, 0),
				MaxAge:   -1,
				HttpOnly: true,
			})
		} else if err != nil {
			log.Printf("Database error: %v", err)
			RenderError(w, r, "Database Error", http.StatusInternalServerError)
			return
		}
	}

	// Handle POST request (create a new post)
	title := r.FormValue("title")
	content := r.FormValue("content")
	categories := r.Form["category"] // Get multiple categories

	// Validate input
	if title == "" || content == "" || len(categories) == 0 {
		RenderError(w, r, "Title, content, and at least one category are required", http.StatusBadRequest)
		return
	}

	// Insert the new post into the database
	result, err := db.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)", userID, title, content)
	if err != nil {
		log.Printf("Error creating post: %v", err)
		RenderError(w, r, "Error creating post", http.StatusInternalServerError)
		return
	}

	// Get the ID of the newly created post
	postID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error retrieving post ID: %v", err)
		RenderError(w, r, "Error retrieving post ID", http.StatusInternalServerError)
		return
	}

	// Insert categories into the database
	for _, category := range categories {
		_, err = db.Exec("INSERT INTO post_categories (post_id, category) VALUES (?, ?)", postID, category)
		if err != nil {
			log.Printf("Error inserting category: %v", err)
			RenderError(w, r, "Error inserting categories", http.StatusInternalServerError)
			return
		}
	}

	// Redirect to the posts page after successful creation
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}
