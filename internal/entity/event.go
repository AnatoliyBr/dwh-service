package entity

type Event struct {
	EventID   int        `json:"event_id"`
	TimeStamp CustomTime `json:"time_stamp"`
	ServiceID int        `json:"service_id"`
}
