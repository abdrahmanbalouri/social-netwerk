package service

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"social-network/app/repository/model"
)

func HandleEventAction(userID string, groupID string, eventID int, action string) (int, error) {
	// Validate action
	if action != "going" && action != "notGoing" {
		return http.StatusBadRequest, errors.New("invalid action: must be 'going' or 'notGoing'")
	}

	// Check if event exists and belongs to the group
	eventIDCheck, eventTimeStr, err := model.GetEventByIDAndGroup(eventID, groupID)
	if err != nil {
		if err == sql.ErrNoRows {
			return http.StatusNotFound, errors.New("event not found")
		}
		return http.StatusInternalServerError, errors.New("database error while fetching event")
	}

	// Parse event time
	eventTime, err := time.Parse(time.RFC3339, eventTimeStr)
	if err != nil {
		return http.StatusInternalServerError, errors.New("invalid event time format")
	}

	// Check if event already passed
	if eventTime.Before(time.Now().UTC()) {
		return http.StatusBadRequest, errors.New("event already passed")
	}

	// Get current user action for this event
	status, err := model.GetUserEventAction(eventIDCheck, userID)

	if err == sql.ErrNoRows {
		// User hasn't responded yet, insert new action
		err = model.InsertEventAction(eventID, userID, action)
		if err != nil {
			return http.StatusInternalServerError, errors.New("failed to save event response")
		}
		return http.StatusOK, nil
	} else if err != nil {
		// Database error
		return http.StatusInternalServerError, errors.New("database error while checking user response")
	}

	// User already has a response
	if action != status {
		// Changing response (going -> notGoing or vice versa)
		err = model.UpdateEventAction(eventID, userID, action)
		if err != nil {
			return http.StatusInternalServerError, errors.New("failed to update event response")
		}
		return http.StatusOK, nil
	}

	// Same action clicked again, remove the response
	err = model.DeleteEventAction(eventID, userID)
	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to remove event response")
	}
	return http.StatusOK, nil
}
