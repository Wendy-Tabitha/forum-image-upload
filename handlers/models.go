package handlers

// User represents a user in the system
type User struct {
	ID       string // Changed from int to string for UUID support
	Email    string
	Username string
	Password string
}

// Post represents a post in the forum
type Post struct {
	ID       int
	UserID   string // Changed from int to string for UUID support
	Title    string
	Content  string
	Category string
}

// Comment represents a comment on a post
type Comment struct {
	ID      int
	PostID  int
	UserID  string // Changed from int to string for UUID support
	Content string
}

// Session represents a user session
type Session struct {
	SessionID string // Unique session identifier
	UserID    string // User ID associated with the session (UUID)
}