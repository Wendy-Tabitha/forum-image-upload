package handlers

import (
	"html/template"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Query to fetch all posts
	rows, err := db.Query("SELECT id, title, content, category FROM posts")
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

	// Render the index page with posts
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, map[string]interface{}{
		"Posts": posts,
	})
}
