package log

import "time"

// Event represents detailed log message.
// It's recommended to instantiate Event with NewEvent function
type Event struct {
	Time    time.Time   `json:"time"`
	Level   Level       `json:"level"`
	Message interface{} `json:"message"`
	Extra   interface{} `json:"extra,omitempty"`
}

// NewEvent creates a new instance of Event
func NewEvent(level Level, message interface{}, extra interface{}) Event {
	return Event{
		Time:    time.Now(),
		Level:   level,
		Message: message,
		Extra:   extra,
	}
}
