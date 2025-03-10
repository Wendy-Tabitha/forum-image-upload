# Image-upload

## Project Overview
This project is designed to create a web forum that allows users to communicate by creating posts and comments. Registered users have the possibility to create a post containing an image as well as text.

## Features
- When viewing the post, users and guests should see the image associated to it.
- in this project you have to handle at least JPEG, PNG and GIF types.
- **User Authentication**: Secure access with user login and registration.
  - **Registration**: Users can register by providing a unique email, username, and password. Passwords are encrypted before storage.
  - **Login**: Users can log in to access the forum. Sessions are managed using cookies with an expiration date.
  - **Session Management**: Each user can have only one active session at a time.

- **Post Management**: Create and view post(s), .
- **Comments**: Registered users can comment on posts, fostering discussion.
- **Likes and Dislikes**: Registered users can like or dislike posts and comments. The number of likes and dislikes is visible to all users.
- **Filtering**: Users can filter posts by categories, created posts, and liked posts.

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
To install this project, follow these steps:
1. Clone the repository: 
   ```bash
   git clone https://learn.zone01kisumu.ke/git/weakinyi/forum-image-upload
2. Navigate to the project directory:
   ```bash
   cd forum-image-upload
   ```
3. Install the required dependencies:
   ```bash
   go get ./...
   ```

## Usage
To run the project with docker, use the following command:
1. Make it executable with this command
```bash
chmod +x script.sh
```
2. Run with this command
```bash
./script.sh
```
- This script will stop and remove any existing container, build the Docker image, and run the container, making it accessible on port 8080.

## Usage without docker
- Run with
``` go
go run .
```

## Testing & Troubleshooting
To run tests, use:
```bash
go test ./...
```
Common issues:
- **Port Conflict**: If you see a "port already in use" error, check for other applications using port 8080.
- **Database Issues**: Verify your database configuration if you encounter connection problems.

## Contributing
We welcome contributions! Please follow these guidelines:
1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Submit a pull request with a clear description of your changes.


## Authors
- antmusumba - [GitHub Profile](https://github.com/antmusumba)
- weakinyi - [GitHub Profile](https://github.com/Wendy-Tabitha)
- Philip38-hub - [GitHub Profile](https://github.com/Philip38-hub)
- hanapiko - [GitHub Profile](https://github.com/hanapiko)



## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.