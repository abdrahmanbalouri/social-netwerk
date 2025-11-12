package service

import (
	"database/sql"
	"errors"
	"net/http"

	"social-network/app/repository/model"
)

func GetEventsService(GrpID string, UserId string) ([]model.Event, int, error) {
	// Get group events
	events, err := model.GetGroupEvents(GrpID, UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusNotFound, errors.New("group not found or no events available")
		}
		return nil, http.StatusInternalServerError, errors.New("failed to retrieve group events")
	}

	return events, http.StatusOK, nil
}