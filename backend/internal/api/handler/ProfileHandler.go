package handlers

import (
	"fmt"
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
    u.nickname, 
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
    ON f.user_id = ? AND f.follower_id = u.id
WHERE u.id = ?;
`

	row := repository.Db.QueryRow(q, currentUserID, targetUserID)

	var user struct {
		ID          string
		Nickname    string
		Email       string
		About       string
		Privacy     string
		Image       string
		Cover       string
		IsFollowing bool
	}

	err = row.Scan(
		&user.ID,
		&user.Nickname,
		&user.Email,
		&user.About,
		&user.Privacy,
		&user.Image,
		&user.Cover,
		&user.IsFollowing,
	)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		fmt.Println("eezdezdezdezdezdezdezdze", err)
		return
	}

	var followers int

	errr := repository.Db.QueryRow(`SELECT COUNT(*) FROM followers WHERE user_id = ?  `, targetUserID).Scan(&followers)
	if errr != nil {
		return
	}

	var following int

	errr = repository.Db.QueryRow(`SELECT COUNT(*) FROM followers WHERE follower_id = ?  `, targetUserID).Scan(&following)
	if errr != nil {
		return
	}

	var IsPending bool
	errr = repository.Db.QueryRow(` 
            SELECT EXISTS(
        SELECT 1 
        FROM follow_requests 
        WHERE user_id = ? 
          AND follower_id = ?
    )`, targetUserID, currentUserID).Scan(&IsPending)
	if errr != nil {
		return
	}

	profileData := map[string]interface{}{
		"id":          user.ID,
		"nickname":    user.Nickname,
		"email":       user.Email,
		"about":       user.About,
		"privacy":     user.Privacy,
		"image":       user.Image,
		"cover":       user.Cover,
		"isFollowing": user.IsFollowing,
		"isPending":   IsPending,
		"following":   following,
		"followers":   followers,
	}



	fmt.Println("<<<<<<", profileData)

	helper.RespondWithJSON(w, http.StatusOK, profileData)
}
