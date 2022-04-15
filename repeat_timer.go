package common

import (
	"sync"
	"time"
)

type RepeatTimer struct {
	Ch chan time.Time

	mtx    sync.Mutex
	name   string
	ticker *time.Ticker
	quit   chan struct{}
	wg     *sync.WaitGroup
	dur    time.Duration
}

func NewRepeatTimer(name string, dur time.Duration) *RepeatTimer {
	var t = &RepeatTimer{
		Ch:     make(chan time.Time),
		ticker: time.NewTicker(dur),
		quit:   make(chan struct{}),
		wg:     new(sync.WaitGroup),
		name:   name,
		dur:    dur,
	}
	t.wg.Add(1)
	go t.fireRoutine(t.ticker)
	return t
}

func (t *RepeatTimer) fireRoutine(ticker *time.Ticker) {
	for {
		select {
		case t_ := <-ticker.C:
			t.Ch <- t_
		case <-t.quit:

			t.wg.Done()
			return
		}
	}
}

func (t *RepeatTimer) Reset() {
	t.Stop()

	t.mtx.Lock()
	defer t.mtx.Unlock()

	t.ticker = time.NewTicker(t.dur)
	t.quit = make(chan struct{})
	t.wg.Add(1)
	go t.fireRoutine(t.ticker)
}

func (t *RepeatTimer) Stop() bool {
	if t == nil {
		return false
	}
	t.mtx.Lock()
	defer t.mtx.Unlock()

	exists := t.ticker != nil
	if exists {
		t.ticker.Stop()
		select {
		case <-t.Ch:

		default:
		}
		close(t.quit)
		t.wg.Wait()
		t.ticker = nil
	}
	return exists
}
