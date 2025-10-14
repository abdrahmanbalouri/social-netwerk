package handlers

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"social-network/internal/helper"
	"social-network/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Nickname  string `json:"nickname"`
	DateBirth int64  `json:"dob"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Image     string `json:"avatar"`
	About     string `json:"aboutMe"`
	Privacy   string `json:"privacy"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	userInformation := User{
		Nickname:  r.FormValue("nickname"),
		FirstName: r.FormValue("firstName"),
		LastName:  r.FormValue("lastName"),
		Email:     r.FormValue("email"),
		Password:  r.FormValue("password"),
		Image:     "",
		About:     r.FormValue("aboutMe"),
		Privacy:   r.FormValue("privacy"),
	}

	file, handler, err := r.FormFile("avatar")
	if err == nil {
		defer file.Close()
		filePath := fmt.Sprintf("../frontend/my-app/public/uploads/%s", handler.Filename)
		dst, _ := os.Create(filePath)
		defer dst.Close()
		io.Copy(dst, file)
		userInformation.Image = handler.Filename
	} else {
		fmt.Println("err", err)
	}

	dobStr := r.FormValue("dob")
	if dobStr != "" {
		unixDob, err := strconv.ParseInt(dobStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid date format for dob", http.StatusBadRequest)
			return
		}
		userInformation.DateBirth = unixDob
	}
	var exists int
	err = repository.Db.QueryRow(
		"SELECT COUNT(*) FROM users WHERE email = ? OR nickname = ?",
		html.EscapeString(userInformation.Email),
		html.EscapeString(userInformation.Nickname),
	).Scan(&exists)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	if exists > 0 {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(userInformation.Email) {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInformation.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Password encryption failed", http.StatusInternalServerError)
		return
	}

	if len(strings.TrimSpace(userInformation.Nickname)) == 0 ||
		len(strings.TrimSpace(userInformation.FirstName)) == 0 ||
		len(strings.TrimSpace(userInformation.LastName)) == 0 ||
		userInformation.Email == "" || userInformation.Password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	age := time.Now().Year() - time.Unix(userInformation.DateBirth, 0).Year()
	if age < 13 || age > 120 {
		http.Error(w, "Invalid age", http.StatusBadRequest)
		return
	}
	id := uuid.New()
	res, err := repository.Db.Exec(`
	INSERT INTO users 
	(id, nickname, date_birth, first_name, last_name, email, password, image, about, privacy, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id,
		html.EscapeString(userInformation.Nickname),
		userInformation.DateBirth,
		html.EscapeString(userInformation.FirstName),
		html.EscapeString(userInformation.LastName),
		html.EscapeString(userInformation.Email),
		hashedPassword,
		func() string {
			if userInformation.Image == "" {
				return "default.png"
			}
			return userInformation.Image
		}(),
		func() string {
			if userInformation.About == "" {
				return "No description"
			}
			return userInformation.About
		}(),
		func() string {
			if userInformation.Privacy == "" {
				return "public"
			}
			return userInformation.Privacy
		}(),
		time.Now().Unix(),
	)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Error executing query")

		return
	}
	_, err = res.LastInsertId()
	if err != nil {

		http.Error(w, "Error retrieving user ID", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"User registered successfully"}`))
}
