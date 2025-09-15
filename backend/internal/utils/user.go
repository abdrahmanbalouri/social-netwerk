package utils

type LoginInformation struct {
	Nickname string `json:"email"`
	Password string `json:"password"`
}

type PostRequest struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	PathImage string `json:"path"`
}
