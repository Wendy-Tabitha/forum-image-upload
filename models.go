package main

type User struct {
    ID       int
    Email    string
    Username string
    Password string
}

type Post struct {
    ID       int
    UserID   int
    Title    string
    Content  string
    Category string
}

type Comment struct {
    ID      int
    PostID  int
    UserID  int
    Content string
}