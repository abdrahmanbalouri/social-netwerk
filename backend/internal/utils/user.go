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
type User struct {
	ID       string `json:"id"`
	Image    string `json:"image"` 
	First_name string `json:"first_name"` 
	Last_name string `json:"last_name"` 

}
