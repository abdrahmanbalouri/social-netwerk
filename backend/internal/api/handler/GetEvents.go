package handlers

import (
	"net/http"
	"strings"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

func GetEvents(w http.ResponseWriter, r *http.Request) {
	type Event struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Date        string `json:"time"`
		Going       int    `json:"going"`
		NotGoing    int    `json:"notGoing"`
		UserAction  string `json:"userAction"`
	}
	UserId, err := helper.AuthenticateUser(r)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		helper.RespondWithError(w, http.StatusNotFound, "Group not found")
		return
	}
	GrpID := parts[3]

	var exists string

	err = repository.Db.QueryRow(`select  user_id from group_members where user_id = ? and group_id = ? `, UserId, GrpID).Scan(&exists)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized "+err.Error())
		return
	}

	var Events []Event
	Fquery := `SELECT e.id, e.title, e.description,e.time,
  SUM(CASE WHEN A.action = 'going' THEN 1 ELSE 0 END) AS going,
  SUM(CASE WHEN A.action = 'notGoing' THEN 1 ELSE 0 END) AS notGoing
FROM events e
LEFT JOIN event_Actions A ON e.id = A.event_id
WHERE e.group_id = ?
GROUP BY e.id 
ORDER BY e.time DESC
`
	rows, err := repository.Db.Query(Fquery, GrpID)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Database query error"+err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {

		var events Event

		if err := rows.Scan(&events.ID, &events.Title, &events.Description, &events.Date, &events.Going, &events.NotGoing); err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Database scan error"+err.Error())
			return
		}
		err = repository.Db.QueryRow(`select  action from event_actions where user_id = ? and event_id = ? `, UserId, events.ID).Scan(&events.UserAction)
		if err != nil {
			events.UserAction = ""
		}
		Events = append(Events, events)
	}

	if err := rows.Err(); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Database rows error")
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, Events)
}
