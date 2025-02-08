package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"
)

// Comment handler for processing form submissions
func CommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	postID := r.FormValue("post_id")
	comment := r.FormValue("content")
	userID := getUserIDFromSession(w, r) // Fetch user ID from session

	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	_, err := db.Exec("INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)", postID, userID, comment)
	if err != nil {
		http.Error(w, "Error posting comment", http.StatusInternalServerError)
		log.Println("Error inserting comment:", err)
		return
	}

	// Redirect back to the post
	http.Redirect(w, r, "/post?id="+postID, http.StatusSeeOther)
}

// Fetch comments for a specific post
func GetCommentsForPost(postID int) ([]Comment, error) {
	rows, err := db.Query(`
		SELECT c.id, c.post_id, CAST(c.user_id AS CHAR), c.content, c.created_at, u.username
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ?
		ORDER BY c.created_at DESC
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		var createdAt time.Time
		if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &createdAt, &c.Username); err != nil {
			return nil, err
		}
		c.CreatedAt = createdAt.Format("Jan 2, 2006 at 3:04 PM")
		comments = append(comments, c)
	}
	return comments, nil
}

// 🔹 Get user ID from session
func getUserIDFromSession(w http.ResponseWriter, r *http.Request) string {
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		return ""
	}

	var userID string
	err = db.QueryRow("SELECT user_id FROM sessions WHERE session_id = ?", sessionCookie.Value).Scan(&userID)
	if err == sql.ErrNoRows {
		// Clear invalid session
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HttpOnly: true,
		})
		http.Error(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
		return ""
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return ""
	}

	return userID
}

// 🔹 Fetch a single post by ID
func GetPostByID(id string) (Post, error) {
	var post Post
	err := db.QueryRow("SELECT id, title, content FROM posts WHERE id = ?", id).Scan(&post.ID, &post.Title, &post.Content)
	return post, err
}
