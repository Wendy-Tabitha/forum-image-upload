package handlers

import (
	"net/http"
)

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		postID := r.FormValue("post_id")
		comment := r.FormValue("comment")
		userID := "USER_ID" // Replace with actual user ID from session or cookie

		// Insert the comment into the database
		_, err := db.Exec("INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)", postID, userID, comment)
		if err != nil {
			http.Error(w, "Error posting comment", http.StatusInternalServerError)
			return
		}

		// Redirect back to the post or home page
		http.Redirect(w, r, "/post", http.StatusSeeOther)
		return
	}
}
