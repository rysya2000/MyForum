package models

type User struct {
	UserID   int
	Username string
	Email    string
	Password string
}

type Post struct {
	PostID   int
	Title    string
	Author   string
	Content  string
	Tags     []string
	Like     int
	Dislike  int
	Created  string
	Comments []Comment
}

type Comment struct {
	CommentID int
	Author    string
	Content   string
	Like      int
	Dislike   int
}

type Tag struct {
	PostID   int
	Hashtags []string
}
