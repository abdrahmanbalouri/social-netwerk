package handlers

import (
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"social-network/app/helper"
	regestire "social-network/app/repository/registre"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

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

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form and files
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Error parsing form")
		return
	}

	// Collect form data
	nickname := r.FormValue("nickname")
	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	email := r.FormValue("email")
	password := r.FormValue("password")
	about := r.FormValue("aboutMe")
	privacy := r.FormValue("privacy")
	dobStr := r.FormValue("dob")

	// Validate required fields
	if strings.TrimSpace(firstName) == "" || strings.TrimSpace(lastName) == "" || email == "" || password == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	// Validate email
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(email) {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid email")
		return
	}

	// Parse DOB and validate age
	var dob int64
	if dobStr != "" {
		dob, err = strconv.ParseInt(dobStr, 10, 64)
		if err != nil {
			helper.RespondWithError(w, http.StatusBadRequest, "Invalid date format")
			return
		}
		age := time.Now().Year() - time.Unix(dob, 0).Year()
		if age < 13 || age > 120 {
			helper.RespondWithError(w, http.StatusBadRequest, "Invalid age")
			return
		}
	}

	// Check if user exists (use registre package)
	exists, err := regestire.UserExists(email, nickname)
	if err != nil {
		log.Printf("UserExists DB error: %v", err)
		helper.RespondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if exists {
		helper.RespondWithError(w, http.StatusConflict, "User already exists")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("bcrypt error: %v", err)
		helper.RespondWithError(w, http.StatusInternalServerError, "Password encryption failed")
		return
	}

	// Handle avatar upload
	var imageFileName string
	file, handler, err := r.FormFile("avatar")
	if err != nil {
		// No file uploaded is acceptable; only fail on real errors
		if err != http.ErrMissingFile {
			log.Printf("FormFile error: %v", err)
			helper.RespondWithError(w, http.StatusBadRequest, "Error retrieving avatar")
			return
		}
	} else {
		defer file.Close()
		uploadDir := "../frontend/public/uploads"
		if mkErr := os.MkdirAll(uploadDir, os.ModePerm); mkErr != nil {
			log.Printf("mkdir error: %v", mkErr)
			helper.RespondWithError(w, http.StatusInternalServerError, "Unable to create upload directory")
			return
		}

		// sanitize filename and generate unique name
		safeFilename := filepath.Base(handler.Filename)
		ext := filepath.Ext(safeFilename)
		newName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		filePath := filepath.Join(uploadDir, newName)

		dst, createErr := os.Create(filePath)
		if createErr != nil {
			log.Printf("create file error: %v", createErr)
			helper.RespondWithError(w, http.StatusInternalServerError, "Unable to save avatar")
			return
		}
		defer dst.Close()

		if _, copyErr := io.Copy(dst, file); copyErr != nil {
			log.Printf("copy file error: %v", copyErr)
			helper.RespondWithError(w, http.StatusInternalServerError, "Unable to save avatar")
			return
		}

		imageFileName = newName
	}

	// Prepare user struct (handler-local)
	user := User{
		ID:        uuid.New().String(),
		Nickname:  html.EscapeString(nickname),
		DateBirth: dob,
		FirstName: html.EscapeString(firstName),
		LastName:  html.EscapeString(lastName),
		Email:     html.EscapeString(email),
		Password:  hashedPassword,
		Image:     imageFileName,
		About:     "No description",
		Privacy:   "public",
		CreatedAt: time.Now().Unix(),
	}
	if about != "" {
		user.About = about
	}
	if privacy != "" {
		user.Privacy = privacy
	}

	// convert handler.User -> repository-local User to avoid import cycle
	regUser := regestire.User{
		ID:        user.ID,
		Nickname:  user.Nickname,
		DateBirth: user.DateBirth,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
		Image:     user.Image,
		About:     user.About,
		Privacy:   user.Privacy,
		CreatedAt: user.CreatedAt,
	}

	// Insert user in DB
	if err := regestire.CreateUser(regUser); err != nil {
		log.Printf("CreateUser error: %v", err)
		helper.RespondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	helper.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "User registered successfully"})
}
