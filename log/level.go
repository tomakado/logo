package log

// Level represents logging level.
type Level struct {
	value uint8
	repr  string
}

// NewLevel creates a new instance of Level with given value and string representation.
func NewLevel(value uint8, repr string) Level {
	return Level{value, repr}
}

// String converts level to string representation.
func (l Level) String() string {
	return l.repr
}

// Gt returns true if numeric representation of level is greater than
// numeric representation of given other level.
func (l Level) Gt(other Level) bool {
	return l.value > other.value
}

// Gte returns true if numeric representation of level is greater than
// of equal to numeric representation of given other level.
func (l Level) Gte(other Level) bool {
	return l.value >= other.value
}

// Supported logging levels.
var (
	LevelVerbose   Level = NewLevel(10, "VERBOSE")
	LevelImportant Level = NewLevel(20, "IMPORTANT")
)
