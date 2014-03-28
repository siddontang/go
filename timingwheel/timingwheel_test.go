package timingwheel

import (
	"testing"
	"time"
)

func TestTimingWheel(t *testing.T) {
	w := NewTimingWheel(1*time.Second, 10)

	println(time.Now().Unix())
	for {
		select {
		case <-w.After(1 * time.Second):
			println(time.Now().Unix())
			return
		}
	}
}

func TestTask(t *testing.T) {
	w := NewTimingWheel(1*time.Second, 10)

	r := make(chan struct{})
	f := func() {
		println("hello world")
		r <- struct{}{}
	}

	w.AddTask(1*time.Second, f)

	<-r
	println("over")
}
