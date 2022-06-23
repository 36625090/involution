package logging

import (
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	t.Log(nextHourOfMilliDuration(), time.Now().Add(nextHourOfMilliDuration()))
	t.Log(nextDayOfMilliDuration(), time.Now().Add(nextDayOfMilliDuration()))
}
