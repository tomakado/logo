package log

import "time"

// Event represents detailed log message.
// It's recommended to instantiate Event with NewEvent function.
type Event struct {
	Time    time.Time   `json:"time"`
	Level   Level       `json:"-"`
	Message interface{} `json:"message"`
	Extra   Extra       `json:"extra,omitempty"`
}

// Extra is a set of key-value pairs (map in other words) used to
// extend the context of log event.
type Extra map[string]interface{}

// NewEvent creates a new instance of Event.
func NewEvent(level Level, message interface{}, extra Extra) Event {
	notNilExtra := extra
	if notNilExtra == nil {
		notNilExtra = Extra{}
	}

	return Event{
		Time:    time.Now(),
		Level:   level,
		Message: message,
		Extra:   notNilExtra,
	}
}
