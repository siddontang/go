package timingwheel

import (
	"github.com/siddontang/golib/log"
	"sync"
	"time"
)

type TaskFunc func()

type bucket struct {
	c     chan struct{}
	tasks []TaskFunc
}

const defaultTasksSize = 16

type TimingWheel struct {
	sync.Mutex

	interval time.Duration

	ticker *time.Ticker
	quit   chan struct{}

	maxTimeout time.Duration

	buckets []bucket

	pos int
}

func NewTimingWheel(interval time.Duration, buckets int) *TimingWheel {
	w := new(TimingWheel)

	w.interval = interval

	w.quit = make(chan struct{})
	w.pos = 0

	w.maxTimeout = time.Duration(interval * (time.Duration(buckets)))

	w.buckets = make([]bucket, buckets)

	for i := range w.buckets {
		w.buckets[i].c = make(chan struct{})
		w.buckets[i].tasks = make([]TaskFunc, 0, defaultTasksSize)
	}

	w.ticker = time.NewTicker(interval)
	go w.run()

	return w
}

func (w *TimingWheel) Stop() {
	close(w.quit)
}

func (w *TimingWheel) After(timeout time.Duration) <-chan struct{} {
	if timeout >= w.maxTimeout {
		panic("timeout too much, over maxtimeout")
	}

	w.Lock()

	index := (w.pos + int(timeout/w.interval)) % len(w.buckets)

	b := w.buckets[index].c

	w.Unlock()

	return b
}

func (w *TimingWheel) AddTask(timeout time.Duration, f TaskFunc) {
	if timeout >= w.maxTimeout {
		panic("timeout too much, over maxtimeout")
	}

	w.Lock()

	index := (w.pos + int(timeout/w.interval)) % len(w.buckets)

	w.buckets[index].tasks = append(w.buckets[index].tasks, f)

	w.Unlock()
}

func (w *TimingWheel) run() {
	for {
		select {
		case <-w.ticker.C:
			w.onTicker()
		case <-w.quit:
			w.ticker.Stop()
			return
		}
	}
}

func (w *TimingWheel) onTicker() {
	w.Lock()

	lastC := w.buckets[w.pos].c
	tasks := w.buckets[w.pos].tasks

	w.buckets[w.pos].c = make(chan struct{})
	w.buckets[w.pos].tasks = w.buckets[w.pos].tasks[0:0:defaultTasksSize]

	w.pos = (w.pos + 1) % len(w.buckets)

	w.Unlock()

	close(lastC)

	if len(tasks) > 0 {
		f := func(tasks []TaskFunc) {
			defer func() {
				if e := recover(); e != nil {
					log.Fatal("run task fatal %v", e)
				}
			}()
			for _, task := range tasks {
				task()
			}
		}

		go f(tasks)
	}
}
