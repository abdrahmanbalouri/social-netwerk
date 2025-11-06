package model

import (
	"social-network/internal/repository"
)

type Event struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"time"`
	Going       int    `json:"going"`
	NotGoing    int    `json:"notGoing"`
	UserAction  string `json:"userAction"`
}
func GetGroupEvents(GrpID string, UserId string) ([]Event, error) {
	//  COUNT(*) FILTER (WHERE action = 'going') AS going // bhalha 
	var Events []Event
	Fquery := `SELECT e.id, e.title, e.description,e.time,
  SUM(CASE WHEN A.action = 'going' THEN 1 ELSE 0 END) AS going,
  SUM(CASE WHEN A.action = 'notGoing' THEN 1 ELSE 0 END) AS notGoing
FROM events e
LEFT JOIN event_Actions A ON e.id = A.event_id
WHERE e.group_id = ?
GROUP BY e.id 
ORDER BY e.time DESC
`
	rows, err := repository.Db.Query(Fquery, GrpID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var events Event

		if err := rows.Scan(&events.ID, &events.Title, &events.Description, &events.Date, &events.Going, &events.NotGoing); err != nil {
			return nil, err
		}
		err = repository.Db.QueryRow(`select  action from event_actions where user_id = ? and event_id = ? `, UserId, events.ID).Scan(&events.UserAction)
		if err != nil {
			events.UserAction = ""
		}
		Events = append(Events, events)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return Events, nil
}
