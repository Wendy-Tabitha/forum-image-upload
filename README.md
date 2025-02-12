# Web Forum Project

## Overview

This project involves creating a web forum that allows users to communicate by creating posts and comments. The forum supports features such as liking/disliking posts and comments, associating categories with posts, and filtering posts based on categories, created posts, and liked posts. The project uses **SQLite** for database management and **Docker** for containerization.

## Features

### User Authentication
- **Registration**: Users can register by providing a unique email, username, and password. Passwords are encrypted before storage.
- **Login**: Users can log in to access the forum. Sessions are managed using cookies with an expiration date.
- **Session Management**: Each user can have only one active session at a time. UUIDs are used for session management (Bonus).

### Communication
- **Posts**: Registered users can create posts and associate them with one or more categories.
- **Comments**: Registered users can comment on posts.
- **Visibility**: Posts and comments are visible to all users (registered and non-registered).

### Likes and Dislikes
- **Likes/Dislikes**: Registered users can like or dislike posts and comments.
- **Count Visibility**: The number of likes and dislikes is visible to all users.

### Filtering
- **Categories**: Users can filter posts by categories.
- **Created Posts**: Registered users can filter posts they have created.
- **Liked Posts**: Registered users can filter posts they have liked.

---

## Technologies Used

- **Backend**: Go (Golang)
- **Database**: SQLite
- **Frontend**: HTML, CSS, JavaScript (no frameworks or libraries)
- **Containerization**: Docker
- **Password Encryption**: bcrypt (Bonus)
- **Session Management**: UUID (Bonus)

---

## Setup Instructions

### Prerequisites
- Docker installed on your machine.
- Basic knowledge of Go and SQL.

### Steps to Run the Project

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/your-username/your-repo-name.git
   cd your-repo-name

2. Build the Docker image
```bash
docker build -t web-forum .
```

3. Run the Docker Container
```bash
docker run -p 8080:8080 web-forum
```
4. Access the Forum
- Open your browser and go to http://localhost:8080.

### Error Handling
- HTTP Status Codes: Proper HTTP status codes are returned for errors (e.g., 400 Bad Request, 401 Unauthorized, 404 Not Found, 500 Internal Server Error).

- User-Friendly Messages: Error responses include user-friendly messages.

### License
This project is licensed under the MIT License. See the LICENSE file for details.