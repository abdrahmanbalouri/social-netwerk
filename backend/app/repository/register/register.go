package regestire

import "social-network/app/repository"

// add a repository-local User model to avoid importing the handler package
type User struct {
	ID        string
	Nickname  string
	DateBirth int64
	FirstName string
	LastName  string
	Email     string
	Password  []byte
	Image     string
	About     string
	Privacy   string
	CreatedAt int64
}

// Check if a user with the given email or nickname exists
func UserExists(email, nickname string) (bool, error) {
	var count int
	err := repository.Db.QueryRow(
		`SELECT COUNT(*) FROM users WHERE (email = ? OR nickname = ?) AND nickname != ""`,
		email, nickname,
	).Scan(&count)
	return count > 0, err
}

// Insert a new user into the database
func CreateUser(user User) error {
	_, err := repository.Db.Exec(`
		INSERT INTO users 
		(id, nickname, date_birth, first_name, last_name, email, password, image, about, privacy, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.ID,
		user.Nickname,
		user.DateBirth,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.Image,
		user.About,
		user.Privacy,
		user.CreatedAt,
	)
	return err
}
