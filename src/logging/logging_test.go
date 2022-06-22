package logging

import (
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {

	next := time.Now().Add(nextDayLeft())

	lastDay := next.Add(-rotationInterval).Format("2006-01-02")
	t.Log(lastDay, int64(rotationInterval))

	t.Log(nextDayLeft(), time.Now().Add(nextDayLeft()))
}
