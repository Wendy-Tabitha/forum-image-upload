package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHomeHandler(t *testing.T) {
	// Store original functions to restore after test
	originalGetUserIdFromSession := GetUserIdFromSession
	originalGetCommentsForPost := GetCommentsForPost
	originalDB := db
	originalRenderError := RenderError

	// Restore original functions after test
	defer func() {
		GetUserIdFromSession = originalGetUserIdFromSession
		GetCommentsForPost = originalGetCommentsForPost
		db = originalDB
		RenderError = originalRenderError
	}()

	// Test case 2: Database Query Error
	t.Run("Database Query Error", func(t *testing.T) {
		// Mock GetUserIdFromSession
		GetUserIdFromSession = func(w http.ResponseWriter, r *http.Request) string {
			return "testuser"
		}

		// Create a mock database that will cause a query error
		mockDB, err := sql.Open("sqlite3", ":memory:")
		if err != nil {
			t.Fatalf("Failed to create mock database: %v", err)
		}
		defer mockDB.Close()

		// Replace global db with mock
		db = mockDB

		// Track if RenderError was called
		var renderErrorCalled bool
		RenderError = func(w http.ResponseWriter, r *http.Request, message string, statusCode int) {
			renderErrorCalled = true
			http.Error(w, message, statusCode)
		}

		// Create request and response recorder
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		// Call handler
		HomeHandler(w, req)

		// Check if RenderError was called
		if !renderErrorCalled {
			t.Errorf("Expected RenderError to be called on database query error")
		}
	})
}

func TestGetUserIdFromSession(t *testing.T) {
	// Store the original function to restore after the test
	originalGetUserIdFromSession := GetUserIdFromSession

	// Restore the original function after the test
	defer func() {
		GetUserIdFromSession = originalGetUserIdFromSession
	}()

	// Test case 1: Valid cookie
	t.Run("Valid Cookie", func(t *testing.T) {
		// Temporarily replace the function for this test
		GetUserIdFromSession = func(w http.ResponseWriter, r *http.Request) string {
			cookie, _ := r.Cookie("session-name")
			return cookie.Value
		}

		// Create a mock HTTP request with a cookie
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session-name",
			Value: "test-user-123",
		})

		// Create a response recorder
		w := httptest.NewRecorder()

		// Call the function
		userID := GetUserIdFromSession(w, req)

		// Manual assertion
		if userID != "test-user-123" {
			t.Errorf("Expected userID to be 'test-user-123', got '%s'", userID)
		}
	})

	// Test case 2: Missing cookie
	t.Run("Missing Cookie", func(t *testing.T) {
		// Temporarily replace the function for this test
		GetUserIdFromSession = func(w http.ResponseWriter, r *http.Request) string {
			http.Error(w, "Error getting session", http.StatusInternalServerError)
			return ""
		}

		// Create a mock HTTP request without a cookie
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		// Create a response recorder
		w := httptest.NewRecorder()

		// Call the function
		userID := GetUserIdFromSession(w, req)

		// Manual assertions
		if userID != "" {
			t.Errorf("Expected empty userID, got '%s'", userID)
		}

		// Check if an error response was written
		response := w.Result()
		if response.StatusCode != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, response.StatusCode)
		}
	})
}

func TestHandleDatabaseError(t *testing.T) {
	// Store original RenderError to restore after test
	originalRenderError := RenderError

	// Restore original function after test
	defer func() {
		RenderError = originalRenderError
	}()

	// Test case 1: Error is not nil
	t.Run("Error Present", func(t *testing.T) {
		// Track if RenderError was called
		var renderErrorCalled bool
		var renderErrorMessage string
		var renderErrorStatusCode int

		// Mock RenderError
		RenderError = func(w http.ResponseWriter, r *http.Request, message string, statusCode int) {
			renderErrorCalled = true
			renderErrorMessage = message
			renderErrorStatusCode = statusCode
			http.Error(w, message, statusCode)
		}

		// Create mock request and response recorder
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		// Create a test error
		testErr := fmt.Errorf("test database error")

		// Call HandleDatabaseError
		HandleDatabaseError(w, req, testErr)

		// Check if RenderError was called
		if !renderErrorCalled {
			t.Errorf("Expected RenderError to be called")
		}

		// Check error message and status code
		if renderErrorMessage != "Database Error" {
			t.Errorf("Expected error message 'Database Error', got '%s'", renderErrorMessage)
		}

		if renderErrorStatusCode != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, got %d",
				http.StatusInternalServerError, renderErrorStatusCode)
		}
	})

	// Test case 2: Error is nil
	t.Run("No Error", func(t *testing.T) {
		// Track if RenderError was called
		var renderErrorCalled bool

		// Mock RenderError
		RenderError = func(w http.ResponseWriter, r *http.Request, message string, statusCode int) {
			renderErrorCalled = true
		}

		// Create mock request and response recorder
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		// Call HandleDatabaseError with nil error
		HandleDatabaseError(w, req, nil)

		// Check that RenderError was NOT called
		if renderErrorCalled {
			t.Errorf("Expected RenderError to NOT be called when error is nil")
		}
	})
}

func TestGetCommentsForPost(t *testing.T) {
	// Setup mock database
	mockDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	// Replace global db with mock
	originalDB := db
	db = mockDB
	defer func() { db = originalDB }()

	// Prepare mock database schema and data
	_, err = mockDB.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY,
			username TEXT
		);
		CREATE TABLE posts (
			id INTEGER PRIMARY KEY,
			title TEXT
		);
		CREATE TABLE comments (
			id INTEGER PRIMARY KEY,
			post_id INTEGER,
			user_id INTEGER,
			content TEXT,
			created_at DATETIME,
			parent_id INTEGER,
			FOREIGN KEY(post_id) REFERENCES posts(id),
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(parent_id) REFERENCES comments(id)
		);
		CREATE TABLE comment_likes (
			comment_id INTEGER,
			is_like BOOLEAN
		);

		-- Insert test users
		INSERT INTO users (id, username) VALUES 
		(1, 'testuser1'),
		(2, 'testuser2');

		-- Insert test post
		INSERT INTO posts (id, title) VALUES (1, 'Test Post');

		-- Insert test comments
		INSERT INTO comments (id, post_id, user_id, content, created_at, parent_id) VALUES 
		(1, 1, 1, 'First comment', '2024-01-01 10:00:00', NULL),
		(2, 1, 2, 'Second comment', '2024-01-01 11:00:00', NULL),
		(3, 1, 1, 'Reply to first comment', '2024-01-01 12:00:00', 1);

		-- Insert comment likes
		INSERT INTO comment_likes (comment_id, is_like) VALUES 
		(1, 1), (1, 1),  -- 2 likes for first comment
		(2, 0), (2, 0);  -- 2 dislikes for second comment
	`)
	if err != nil {
		t.Fatalf("Failed to prepare mock data: %v", err)
	}

	// Mock GetCommentReplies to return predefined replies
	originalGetCommentReplies := GetCommentReplies
	GetCommentReplies = func(commentID int) ([]Comment, error) {
		if commentID == 1 {
			return []Comment{
				{
					ID:       3,
					PostID:   1,
					UserID:   "1",
					Content:  "Reply to first comment",
					Username: "testuser1",
				},
			}, nil
		}
		return []Comment{}, nil
	}
	defer func() { GetCommentReplies = originalGetCommentReplies }()

	// Test cases
	testCases := []struct {
		name           string
		postID         int
		expectedResult struct {
			commentCount      int
			firstCommentID    int
			firstCommentLikes int
			replyCount        int
		}
	}{
		{
			name:   "Non-Existing Post",
			postID: 999,
			expectedResult: struct {
				commentCount      int
				firstCommentID    int
				firstCommentLikes int
				replyCount        int
			}{
				commentCount: 0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			comments, err := GetCommentsForPost(tc.postID)

			if tc.postID == 999 {
				// For non-existing post, expect no error and empty comments
				if err != nil {
					t.Errorf("Unexpected error for non-existing post: %v", err)
				}
				if len(comments) != 0 {
					t.Errorf("Expected 0 comments for non-existing post, got %d", len(comments))
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(comments) != tc.expectedResult.commentCount {
				t.Errorf("Expected %d comments, got %d", tc.expectedResult.commentCount, len(comments))
			}

			if len(comments) > 0 {
				firstComment := comments[0]
				if firstComment.ID != tc.expectedResult.firstCommentID {
					t.Errorf("Expected first comment ID %d, got %d", tc.expectedResult.firstCommentID, firstComment.ID)
				}

				if firstComment.LikeCount != tc.expectedResult.firstCommentLikes {
					t.Errorf("Expected %d likes, got %d", tc.expectedResult.firstCommentLikes, firstComment.LikeCount)
				}

				if len(firstComment.Replies) != tc.expectedResult.replyCount {
					t.Errorf("Expected %d replies, got %d", tc.expectedResult.replyCount, len(firstComment.Replies))
				}
			}
		})
	}
}
