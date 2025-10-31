package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		helper.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	userID, err := helper.AuthenticateUser(r)
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

	err = repository.Db.QueryRow(`select  user_id from group_members where user_id = ? and group_id = ? `, userID, GrpID).Scan(&exists)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized " + err.Error() )
		return
	}

jsonDecoder := json.NewDecoder(r.Body)
	var event struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		DateTime string `json:"dateTime"`

	}
	err = jsonDecoder.Decode(&event)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	result, err := repository.Db.Exec(`INSERT INTO events (group_id, title, description, time) VALUES (?, ?, ?,?)`, GrpID, event.Title, event.Description  , event.DateTime)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Database insert error" + err.Error())
		return
	}

	eventID, err := result.LastInsertId()
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Database retrieval error")
		return
	}

	type CreatedEvent struct {
		ID          int64  `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`	
		Date  string `json:"time"`				



	}
	

	helper.RespondWithJSON(w, http.StatusOK, CreatedEvent{
		ID:          eventID,
		Title:       event.Title,
		Description: event.Description,
		Date: event.DateTime,
	})

}
