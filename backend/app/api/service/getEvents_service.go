package service

import "social-network/app/repository/model"

func GetEventsService(GrpID string, UserId string) ([]model.Event, error) {
	return model.GetGroupEvents(GrpID, UserId)
}
