package log_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tomakado/logo/log"
)

func TestLevel_Gt(t *testing.T) {
	assert.True(t, log.LevelImportant.Gt(log.LevelVerbose))
	assert.False(t, log.LevelVerbose.Gt(log.LevelImportant))
}

func TestLevel_Gte(t *testing.T) {
	assert.True(t, log.LevelImportant.Gte(log.LevelVerbose))
	assert.True(t, log.LevelImportant.Gte(log.LevelImportant))
	assert.False(t, log.LevelVerbose.Gte(log.LevelImportant))
}
