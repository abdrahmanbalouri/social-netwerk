package handlers

import (
	"net/http"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	targetUserID := r.URL.Query().Get("userId")
	currentUserID, err := helper.AuthenticateUser(r)
	if err != nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if targetUserID == "0" || targetUserID == "" {
		targetUserID = currentUserID
	}

	q := `
SELECT 
    u.id, 
    u.first_name, u.last_name, 
    u.email, 
    u.about, 
    u.privacy, 
    u.image, 
    u.cover,
    CASE 
        WHEN f.follower_id IS NOT NULL THEN 1 
        ELSE 0 
    END AS is_following
FROM users u
LEFT JOIN followers f 
    ON f.user_id = u.id AND f.follower_id = ?
WHERE u.id = ?;
`

	row := repository.Db.QueryRow(q, currentUserID, targetUserID)

	var user struct {
		ID          string
		first_name  string
		last_name   string
		Email       string
		About       string
		Privacy     string
		Image       string
		Cover       string
		IsFollowing bool
	}

	err = row.Scan(
		&user.ID,
		&user.first_name,
		&user.last_name,
		&user.Email,
		&user.About,
		&user.Privacy,
		&user.Image,
		&user.Cover,
		&user.IsFollowing,
	)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var followers, following int
	if err = repository.Db.QueryRow(`SELECT COUNT(*) FROM followers WHERE user_id = ?`, targetUserID).Scan(&followers); err != nil {
		http.Error(w, "Failed to count followers", http.StatusInternalServerError)
		return
	}
	if err = repository.Db.QueryRow(`SELECT COUNT(*) FROM followers WHERE follower_id = ?`, targetUserID).Scan(&following); err != nil {
		http.Error(w, "Failed to count following", http.StatusInternalServerError)
		return
	}

	var isPending bool
	if err = repository.Db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM follow_requests 
			WHERE user_id = ? AND follower_id = ?
		)`, targetUserID, currentUserID).Scan(&isPending); err != nil {
		http.Error(w, "Failed to check pending status", http.StatusInternalServerError)
		return
	}

	profileData := map[string]interface{}{
		"id":          user.ID,
		"first_name":  user.first_name,
		"last_name":   user.last_name,
		"email":       user.Email,
		"about":       user.About,
		"privacy":     user.Privacy,
		"image":       user.Image,
		"cover":       user.Cover,
		"isFollowing": user.IsFollowing,
		"isPending":   isPending,
		"following":   following,
		"followers":   followers,
	}
	helper.RespondWithJSON(w, http.StatusOK, profileData)
}
