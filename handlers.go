package main

import (
	"database/sql"
	"html/template"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
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
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, map[string]interface{}{
		"Posts": posts,
	})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		var user User
		err := db.QueryRow("SELECT id, password FROM users WHERE email = ?", email).Scan(&user.ID, &user.Password)
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// Set session cookie
		http.SetCookie(w, &http.Cookie{
			Name:    "session_id",
			Value:   "some_unique_session_id",       // You should generate a unique session ID
			Expires: time.Now().Add(24 * time.Hour), // Set expiration
		})

		http.Redirect(w, r, "/post", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	tmpl.Execute(w, nil)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")

		// Check if the passwords match
		if password != confirmPassword {
			http.Error(w, "Passwords do not match", http.StatusBadRequest)
			return
		}

		// Check if email already exists
		var existingEmail string
		err := db.QueryRow("SELECT email FROM users WHERE email = ?", email).Scan(&existingEmail)
		if err == nil {
			http.Error(w, "Email already taken", http.StatusBadRequest)
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		// Insert new user into the database
		_, err = db.Exec("INSERT INTO users (email, username, password) VALUES (?, ?, ?)", email, username, hashedPassword)
		if err != nil {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/post", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/register.html"))
	tmpl.Execute(w, nil)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")
		category := r.FormValue("category")

		// Get user ID from session (you should implement session management)
		userID := 1 // Replace with actual user ID from session

		// Insert the new post into the database
		_, err := db.Exec("INSERT INTO posts (user_id, title, content, category) VALUES (?, ?, ?, ?)", userID, title, content, category)
		if err != nil {
			http.Error(w, "Error creating post", http.StatusInternalServerError)
			return
		}
	}

	// Query to fetch all posts (same as indexHandler)
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

	// Render the post page with posts
	tmpl := template.Must(template.ParseFiles("templates/post.html"))
	tmpl.Execute(w, map[string]interface{}{
		"Posts": posts,
	})
}

func likeHandler(w http.ResponseWriter, r *http.Request) {
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

// filterHandler handles filtering posts based on criteria
func filterHandler(w http.ResponseWriter, r *http.Request) {
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
