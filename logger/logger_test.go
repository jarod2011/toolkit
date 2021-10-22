package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLevel_String(t *testing.T) {
	for lvl, want := range map[Level]string{
		Debug: "debug",
		Info:  "info",
		Warn:  "warning",
		Error: "error",
		Fatal: "fatal",
	} {
		assert.Equal(t, lvl.String(), want)
	}
}
