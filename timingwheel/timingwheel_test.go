package timingwheel

import (
	"testing"
	"time"
)

func TestTimingWheel(t *testing.T) {
	w := NewTimingWheel(1*time.Second, 10)

	t.Log(time.Now().Unix())
	for {
		select {
		case <-w.After(1 * time.Second):
			t.Log(time.Now().Unix())
			return
		}
	}
}
