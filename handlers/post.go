package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	var userID string
	sessionCookie, err := r.Cookie("session_id")
	if err == nil {
		err = db.QueryRow("SELECT user_id FROM sessions WHERE session_id = ?", sessionCookie.Value).Scan(&userID)
		if err == sql.ErrNoRows {
			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    "",
				Path:     "/",
				Expires:  time.Unix(0, 0), // Expire immediately
				MaxAge:   -1,
				HttpOnly: true,
			})
		} else if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			log.Println("Error retrieving session:", err)
			return
		}
	}

	if r.Method == http.MethodPost {
		// Handle post creation
		title := r.FormValue("title")
		content := r.FormValue("content")
		categories := r.Form["category"]

		// Validate inputs
		if title == "" || content == "" {
			http.Error(w, "Title and content cannot be empty", http.StatusBadRequest)
			return
		}

		// Insert the new post into the database
		result, err := db.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)", userID, title, content)
		if err != nil {
			http.Error(w, "Error creating post", http.StatusInternalServerError)
			log.Println("Error inserting post:", err)
			return
		}

		// Get the ID of the newly created post
		postID, err := result.LastInsertId()
		if err != nil {
			http.Error(w, "Error getting post ID", http.StatusInternalServerError)
			log.Println("Error getting last insert ID:", err)
			return
		}

		// Insert categories for the post
		for _, category := range categories {
			_, err := db.Exec("INSERT INTO post_categories (post_id, category) VALUES (?, ?)", postID, category)
			if err != nil {
				log.Printf("Error inserting category %s for post %d: %v", category, postID, err)
				// Continue with other categories even if one fails
			}
		}

		// Redirect to home page
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	} else {
		// For GET requests, redirect to home page
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}