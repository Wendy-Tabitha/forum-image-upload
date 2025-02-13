package main

import (
	"log"
	"net/http"

	"github.com/Wendy-Tabitha/discussion/handlers"
)

func main() {
	// Serve static files from the "static" directory
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Define routes
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/post", handlers.PostHandler)
	http.HandleFunc("/like", handlers.LikeHandler)
	http.HandleFunc("/filter", handlers.FilterHandler)
	http.HandleFunc("/comment", handlers.CommentHandler)
	http.HandleFunc("/comment/like", handlers.CommentLikeHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/profile", handlers.ProfileHandler) // Register the profile route

	// Initialize the database
	handlers.InitDB()

	// Start the server
	log.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// func handler(w http.ResponseWriter, r *http.Request) {
// 	switch r.URL.Path {
// 	case "/":
// 		handlers.HomeHandler(w, r)
// 	case "/login":
// 		handlers.LoginHandler(w, r)
// 	case "/register":
// 		handlers.RegisterHandler(w, r)
// 	case "/post":
// 		handlers.PostHandler(w, r)
// 	case "/like":
// 		handlers.LikeHandler(w, r)
// 	case "/filter":
// 		handlers.FilterHandler(w, r)
// 	case "/comment":
// 		handlers.CommentHandler(w, r)
// 	case "/comment/like":
// 		handlers.CommentLikeHandler(w, r)
// 	case "/logout":
// 		handlers.LogoutHandler(w, r)
// 	case "/profile":
// 		handlers.ProfileHandler(w, r)
// 	default:
// 		handlers.RenderError(w, r, "Page not found", http.StatusNotFound)
// 	}
// }