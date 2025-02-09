package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

// var db *sql.DB // Ensure this is properly initialized elsewhere in your code

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Retrieve session cookie
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "You must be logged in to like or dislike a post", http.StatusUnauthorized)
		return
	}

	// Get user ID from session
	var userID string
	err = db.QueryRow("SELECT user_id FROM sessions WHERE session_id = ?", sessionCookie.Value).Scan(&userID)
	if err == sql.ErrNoRows {
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HttpOnly: true,
		})
		http.Error(w, "You must be logged in to like or dislike a post", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Read form values
	postID := r.FormValue("post_id")
	if postID == "" {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	isLike := r.FormValue("is_like") == "true"

	// Check if the user has already liked/disliked this post
	var existingLikeID int
	var existingIsLike bool
	err = db.QueryRow("SELECT id, is_like FROM likes WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&existingLikeID, &existingIsLike)
	if err == sql.ErrNoRows {
		// Insert new like/dislike
		_, err = db.Exec("INSERT INTO likes (user_id, post_id, is_like) VALUES (?, ?, ?)", userID, postID, isLike)
		if err != nil {
			http.Error(w, "Error liking post", http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		http.Error(w, "Error checking like status", http.StatusInternalServerError)
		return
	} else {
		if existingIsLike != isLike {
			_, err = db.Exec("UPDATE likes SET is_like = ? WHERE id = ?", isLike, existingLikeID)
		} else {
			_, err = db.Exec("DELETE FROM likes WHERE id = ?", existingLikeID)
		}
		if err != nil {
			http.Error(w, "Error updating like status", http.StatusInternalServerError)
			return
		}
	}

	// Fetch updated like and dislike counts
	var likeCount, dislikeCount int
	if err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ? AND is_like = true", postID).Scan(&likeCount); err != nil {
		http.Error(w, "Error fetching like count", http.StatusInternalServerError)
		return
	}
	if err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ? AND is_like = false", postID).Scan(&dislikeCount); err != nil {
		http.Error(w, "Error fetching dislike count", http.StatusInternalServerError)
		return
	}

	// Return updated counts as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":       true,
		"like_count":    likeCount,
		"dislike_count": dislikeCount,
	})
}
