package profile

import (
	"net/http"

	"social-network/pkg/db/sqlite"
)

func GetProfileData(targetUserID, currentUserID string, w http.ResponseWriter) (map[string]interface{}, error) {
	q := `
SELECT 
    u.id, 
    u.first_name, u.last_name, 
    u.email, 
    u.about, 
	u.date_birth,
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

	row := sqlite.Db.QueryRow(q, currentUserID, targetUserID)

	var user struct {
		ID          string
		first_name  string
		last_name   string
		Email       string
		About       string
		dateOfBirth int
		Privacy     string
		Image       string
		Cover       string
		IsFollowing bool
	}

	err := row.Scan(
		&user.ID,
		&user.first_name,
		&user.last_name,
		&user.Email,
		&user.About,
		&user.dateOfBirth,
		&user.Privacy,
		&user.Image,
		&user.Cover,
		&user.IsFollowing,
	)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return nil, err
	}

	var followers, following int
	if err = sqlite.Db.QueryRow(`SELECT COUNT(*) FROM followers WHERE user_id = ?`, targetUserID).Scan(&followers); err != nil {
		http.Error(w, "Failed to count followers", http.StatusInternalServerError)
		return nil, err
	}
	if err = sqlite.Db.QueryRow(`SELECT COUNT(*) FROM followers WHERE follower_id = ?`, targetUserID).Scan(&following); err != nil {
		http.Error(w, "Failed to count following", http.StatusInternalServerError)
		return nil, err
	}

	var isPending bool
	if err = sqlite.Db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM follow_requests 
			WHERE user_id = ? AND follower_id = ?
		)`, targetUserID, currentUserID).Scan(&isPending); err != nil {
		http.Error(w, "Failed to check pending status", http.StatusInternalServerError)
		return nil, err
	}

	profileData := map[string]interface{}{
		"id":          user.ID,
		"first_name":  user.first_name,
		"last_name":   user.last_name,
		"email":       user.Email,
		"about":       user.About,
		"privacy":     user.Privacy,
		"dateOfBirth": user.dateOfBirth,
		"image":       user.Image,
		"cover":       user.Cover,
		"isFollowing": user.IsFollowing,
		"isPending":   isPending,
		"following":   following,
		"followers":   followers,
	}
	return profileData, nil
}
