package service

import (
	"database/sql"
	"errors"
	"time"

	"social-network/app/repository/model"
)

func HandleEventAction(userID string, groupID string, eventID int, action string) error {
	if action != "going" && action != "notGoing" {
		return errors.New("invalid action")
	}

	eventIDCheck, eventTimeStr, err := model.GetEventByIDAndGroup(eventID, groupID)
	if err != nil {
		return errors.New("event not found or unauthorized")
	}

	eventTime, err := time.Parse(time.RFC3339, eventTimeStr)
	if err != nil {
		return errors.New("invalid event time format")
	}

	if eventTime.Before(time.Now().UTC()) {
		return errors.New("event already passed")
	}

	status, err := model.GetUserEventAction(eventIDCheck, userID)

	if err == sql.ErrNoRows {
		return model.InsertEventAction(eventID, userID, action)
	} else if err != nil {
		return err
	}

	if action != status {
		return model.UpdateEventAction(eventID, userID, action)
	}

	return model.DeleteEventAction(eventID, userID)
}
