package handlers

import (
	"html/template"
	"net/http"
)

// filterHandler handles filtering posts based on criteria
func FilterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		category := r.URL.Query().Get("category")

		// Query to filter posts by category
		rows, err := db.Query("SELECT id, title, content, category FROM posts WHERE category = ?", category)
		if err != nil {
			http.Error(w, "Error fetching posts", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var posts []Post
		for rows.Next() {
			var post Post
			if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category); err != nil {
				http.Error(w, "Error scanning posts", http.StatusInternalServerError)
				return
			}
			posts = append(posts, post)
		}

		// Render the filtered posts
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.Execute(w, map[string]interface{}{
			"Posts": posts,
		})
		return
	}

	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}
