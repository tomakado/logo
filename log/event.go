package log

import "time"

type Event struct {
	Time    time.Time   `json:"time"`
	Level   Level       `json:"level"`
	Message interface{} `json:"message"`
	Extra   interface{} `json:"extra,omitempty"`
}

func NewEvent(level Level, message interface{}, extra interface{}) Event {
	return Event{
		Time:    time.Now(),
		Level:   level,
		Message: message,
		Extra:   extra,
	}
}
