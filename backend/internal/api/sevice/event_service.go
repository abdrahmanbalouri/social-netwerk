package service

import (
	"errors"
	"social-network/internal/repository"
	"social-network/internal/repository/model"
	"time"
)

func CreateGroupEvent(userID, groupID, title, description, dateTimeStr string) (int64, error) {
	// Validate user in group
	inGroup, err := model.CheckUserInGroup(repository.Db, userID, groupID)
	if err != nil {
		return 0, err
	}
	if !inGroup {
		return 0, errors.New("user is not a member of the group")
	}

	// Validate and parse date
	layout := "2006-01-02T15:04"
	eventTime, err := time.Parse(layout, dateTimeStr)
	if err != nil {
		return 0, errors.New("invalid date format")
	}

	if eventTime.Before(time.Now()) {
		return 0, errors.New("event date and time must be in the future")
	}

	event := model.Event{
		GroupID:     groupID,
		Title:       title,
		Description: description,
		Time:        eventTime,
	}

	// Insert into DB
	eventID, err := model.InsertEvent(repository.Db, event)
	if err != nil {
		return 0, err
	}

	return eventID, nil
}
