package timingwheel

import (
	"time"
)

type regItem struct {
	timeout time.Duration
	reply   chan chan bool
}

type TimingWheel struct {
	interval time.Duration

	ticker *time.Ticker
	quit   chan bool

	reg chan *regItem

	maxTimeout time.Duration
	buckets    []chan bool
	pos        int
}

func NewTimingWheel(interval time.Duration, buckets int) *TimingWheel {
	w := new(TimingWheel)

	w.interval = interval

	w.reg = make(chan *regItem, 128)

	w.quit = make(chan bool)
	w.pos = 0

	w.maxTimeout = time.Duration(interval * (time.Duration(buckets)))

	w.buckets = make([]chan bool, buckets)

	for i := range w.buckets {
		w.buckets[i] = make(chan bool)
	}

	w.ticker = time.NewTicker(interval)
	go w.run()

	return w
}

func (w *TimingWheel) Stop() {
	w.quit <- true
}

func (w *TimingWheel) After(timeout time.Duration) <-chan bool {
	if timeout >= w.maxTimeout {
		panic("timeout too much, over maxtimeout")
	}

	reply := make(chan chan bool)

	w.reg <- &regItem{timeout: timeout, reply: reply}

	return <-reply
}

func (w *TimingWheel) run() {
	for {
		select {
		case item := <-w.reg:
			w.register(item)
		case <-w.ticker.C:
			w.onTicker()
		case <-w.quit:
			w.ticker.Stop()
			return
		}
	}
}

func (w *TimingWheel) register(item *regItem) {
	timeout := item.timeout

	index := (w.pos + int(timeout/w.interval)) % len(w.buckets)

	b := w.buckets[index]

	item.reply <- b
}

func (w *TimingWheel) onTicker() {
	close(w.buckets[w.pos])

	w.buckets[w.pos] = make(chan bool)

	w.pos = (w.pos + 1) % len(w.buckets)
}
