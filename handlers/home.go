package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Check if db is initialized
	if db == nil {
		http.Error(w, "Database connection is not initialized", http.StatusInternalServerError)
		log.Println("Error: Database connection is nil")
		return
	}

	// SQL Query
	query := `
		SELECT 
			p.id, 
			p.title, 
			p.content, 
			COALESCE(GROUP_CONCAT(pc.category), '') AS categories, 
			u.username, 
			p.created_at 
		FROM posts p 
		JOIN users u ON p.user_id = u.id 
		LEFT JOIN post_categories pc ON p.id = pc.post_id 
		GROUP BY p.id, p.title, p.content, u.username, p.created_at  -- Ensure correct GROUP BY clause
		ORDER BY p.created_at DESC
	`

	// Run Query
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Error fetchinccg posts", http.StatusInternalServerError)
		log.Printf("Database query error: %v\nQuery: %s", err, query) // LOG the actual error
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Categories, &post.Username, &post.CreatedAt); err != nil {
			http.Error(w, "Error scanning posts", http.StatusInternalServerError)
			log.Println("Error scanning rows:", err) // Log scanning error
			return
		}
		posts = append(posts, post)
	}

	// Render the index page with posts
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "Error parsing file", http.StatusInternalServerError)
		log.Println("Template parsing error:", err) // Log template error
		return
	}
	if err := tmpl.Execute(w, map[string]interface{}{"Posts": posts}); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Println("Template execution error:", err) // Log rendering error
	}
}
