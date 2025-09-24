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
type User struct{
	Nickname string  `json:"nickname"`
	ID  string     `json:"id"`
      
}
