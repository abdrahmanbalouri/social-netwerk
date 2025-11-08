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
type Group struct {
	ID          string
	Title       string
	Description string
	MemberCount int
}
type GroupRequest struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	InvitedUsers []string `json:"invitedUsers"`
}

type GroupInvitation struct {
	// GroupID      string   `json:"groupID"`
	InvitationType string   `json:"InvitationType"`
	InvitedUsers   []string `json:"invitedUsers"`
}

type GroupResponse struct {
	InvitationType string `json:"invitation_type"`
	InvitationID   string `json:"invitation_id"`
	Response       string `json:"response"`
}

type JoinRequest struct {
	InvitationID string `json:"invitation_id"`
	UserID       string `json:"user_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
}

type FetchGroupInvitation struct {
	GroupID      string `json:"group_id"`
	Title        string `json:"title"`
	InvitationID string `json:"invitation_id"`
}

type Friend struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Image     string `json:"image"`
}