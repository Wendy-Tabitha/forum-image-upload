package handlers

import (
	"database/sql"
	"html/template"
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

		postID, err := result.LastInsertId()
		if err != nil {
			RenderError(w, r, "Error retrieving post ID", http.StatusInternalServerError)
			log.Println("Error retrieving last insert ID:", err)
			return
		}

		for _, category := range categories {
			_, err = db.Exec("INSERT INTO post_categories (post_id, category) VALUES (?, ?)", postID, category)
			if err != nil {
				RenderError(w, r, "Error inserting categories", http.StatusInternalServerError)
				log.Println("Error inserting categories:", err)
				return
			}
		}

		// Redirect back to home page after creating the post
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Query to fetch all posts along with the user's name, creation time, and like/dislike counts
	rows, err := db.Query(`
		SELECT 
			p.id, 
			p.title, 
			p.content,
			COALESCE(GROUP_CONCAT(pc.category), '') as categories,
			u.username, 
			p.created_at,
			COALESCE(SUM(CASE WHEN l.is_like = 1 THEN 1 ELSE 0 END), 0) AS like_count,
			COALESCE(SUM(CASE WHEN l.is_like = 0 THEN 1 ELSE 0 END), 0) AS dislike_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN post_categories pc ON p.id = pc.post_id
		LEFT JOIN likes l ON p.id = l.post_id
		GROUP BY p.id, p.title, p.content, u.username, p.created_at
		ORDER BY p.created_at DESC
	`)
	if err != nil {
		log.Println("Error fetching posts:", err)
		http.Error(w, "Error fetching posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var categories sql.NullString
		if err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&categories,
			&post.Username,
			&post.CreatedAt,
			&post.LikeCount,
			&post.DislikeCount,
		); err != nil {
			RenderError(w, r, "Error scanning posts", http.StatusInternalServerError)
			log.Println("Error scanning post data:", err)
			return
		}
		if categories.Valid {
			post.Categories = categories.String
		} else {
			post.Categories = ""
		}
		posts = append(posts, post)
	}

	// Render the home page with posts
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "Error parsing file", http.StatusInternalServerError)
		log.Println("Error parsing template:", err)
		return
	}

	if err := tmpl.Execute(w, map[string]interface{}{
		"Posts":      posts,
		"IsLoggedIn": userID != "",
	}); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Println("Error executing template:", err)
	}
}
