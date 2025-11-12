package service

import (
    "database/sql"
    "errors"
    "net/http"
    "time"

    "social-network/app/repository/model"
    "social-network/pkg/db/sqlite"
)

func CreateGroupEvent(userID, groupID, title, description, dateTimeStr string) (int64, int, error) {
    // Validate user in group
    inGroup, err := model.CheckUserInGroup(sqlite.Db, userID, groupID)
    if err != nil {
        if err == sql.ErrNoRows {
            return 0, http.StatusNotFound, errors.New("group not found")
        }
        return 0, http.StatusInternalServerError, errors.New("failed to verify group membership")
    }
    if !inGroup {
        return 0, http.StatusForbidden, errors.New("user is not a member of the group")
    }

    // Validate and parse date
    layout := "2006-01-02T15:04"
    eventTime, err := time.Parse(layout, dateTimeStr)
    if err != nil {
        return 0, http.StatusBadRequest, errors.New("invalid date format, expected YYYY-MM-DDTHH:MM")
    }

    if eventTime.Before(time.Now()) {
        return 0, http.StatusBadRequest, errors.New("event date and time must be in the future")
    }

    event := model.Eventt{
        GroupID:     groupID,
        Title:       title,
        Description: description,
        Time:        eventTime,
    }

    // Insert into DB
    eventID, err := model.InsertEvent(sqlite.Db, event)
    if err != nil {
        return 0, http.StatusInternalServerError, errors.New("failed to create event")
    }

    return eventID, http.StatusCreated, nil
}