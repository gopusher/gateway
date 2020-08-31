package log

import (
	"testing"
)

func TestLogMessages(t *testing.T) {
	New(true)
	defer Sync()
	Info("this is a test message")
}
