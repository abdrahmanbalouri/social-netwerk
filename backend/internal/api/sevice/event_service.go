package service

import (
	"errors"
	"time"

	"social-network/internal/repository"
)

func CreateGroupEvent(userID, groupID, title, description, dateTimeStr string) (int64, error) {
	// Validate user in group
	inGroup, err := repository.CheckUserInGroup(repository.Db, userID, groupID)
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

	event := repository.Event{
		GroupID:     groupID,
		Title:       title,
		Description: description,
		Time:        eventTime,
	}

	// Insert into DB
	eventID, err := repository.InsertEvent(repository.Db, event)
	if err != nil {
		return 0, err
	}

	return eventID, nil
}
