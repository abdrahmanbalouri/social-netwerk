package utils

import "time"

type LoginInformation struct {
	Nickname string `json:"email"`
	Password string `json:"password"`
}

type PostRequest struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	PathImage string `json:"path"`
}
type User struct {
	ID         string `json:"id"`
	Image      string `json:"image"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
}
type GroupPost struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	ImagePath     string    `json:"image_path"`
	CreatedAt     time.Time `json:"created_at"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Profile       string    `json:"profile"`
	LikeCount     int       `json:"like"`
	LikedByUser   bool      `json:"liked_by_user"`
	CommentsCount int       `json:"comments_count"`
}
