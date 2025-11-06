package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"social-network/internal/helper"
	"social-network/internal/repository"
)

func EventAction(w http.ResponseWriter, r *http.Request) {
	UserID, err := helper.AuthenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		helper.RespondWithError(w, http.StatusNotFound, "Group not found")
		return
	}
	GrpID := parts[4]

	var exists string

	err = repository.Db.QueryRow(`select  user_id from group_members where user_id = ? and group_id = ? `, UserID, GrpID).Scan(&exists)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized "+err.Error())
		return
	}

	var req struct {
		Action  string `json:"status"`
		EventID int    `json:"eventID"`
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Action != "going" && req.Action != "notGoing" {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return

	}
	var event int
	var timee string
	err = repository.Db.QueryRow(`select  id , time  from events where id = ? and group_id = ? `, req.EventID, GrpID).Scan(&event, &timee)
	if err != nil {
		helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized "+err.Error())
		return
	}

	eventTime, err := time.Parse(time.RFC3339, timee)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid date format")
		return
	}

	currentTime := time.Now().UTC().Add(time.Hour.Abs())

	if (eventTime).Before(currentTime) {
		helper.RespondWithError(w, http.StatusBadRequest, "Event date and time must be in the future 222 ")
		return
	}

	var status string
	err = repository.Db.QueryRow(`select  action    from event_actions where event_id = ? and  user_id= ? `, req.EventID, UserID).Scan(&status)

	if err == sql.ErrNoRows {
		_, err = repository.Db.Exec(`insert into event_actions (event_id,user_id,action) values (?,?,?)`, req.EventID, UserID, req.Action)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "database error  "+err.Error())
			return
		}
	} else if req.Action != status {
		_, err = repository.Db.Exec(`update  event_actions set action = ? where event_id = ? and user_id = ? `, req.Action, req.EventID, UserID)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "database error  "+err.Error())
			return
		}
	} else if req.Action == status {
		_, err = repository.Db.Exec(`delete from event_actions  where event_id = ? and user_id = ? `, req.EventID, UserID)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "database error  "+err.Error())
			return
		}
	} else {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid action")
		return
	}

	var eventReq struct {
		EventID int    `json:"eventID"`
		Status  string `json:"status"`
	}

	helper.RespondWithJSON(w, http.StatusOK, eventReq)
}
