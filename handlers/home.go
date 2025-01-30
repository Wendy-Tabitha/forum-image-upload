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
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "Error parsing file", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, map[string]interface{}{
		"Posts": posts,
	})
}
