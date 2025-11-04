package service

import "social-network/internal/repository/model"

func GetEventsService(GrpID string, UserId string) ([]model.Event, error) {
	return model.GetGroupEvents(GrpID, UserId)
}	