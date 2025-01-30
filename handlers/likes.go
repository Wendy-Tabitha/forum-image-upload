package handlers

import (
	"database/sql"
	"net/http"
)

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		postID := r.FormValue("post_id")
		userID := 1 // Replace with actual user ID from session
		isLike := r.FormValue("is_like") == "true"

		// Check if the user has already liked/disliked this post
		var existingLikeID int
		err := db.QueryRow("SELECT id FROM likes WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&existingLikeID)

		if err == sql.ErrNoRows {
			// Insert new like/dislike
			_, err = db.Exec("INSERT INTO likes (user_id, post_id, is_like) VALUES (?, ?, ?)", userID, postID, isLike)
			if err != nil {
				http.Error(w, "Error liking post", http.StatusInternalServerError)
				return
			}
		}
	}
}
