package internal

import (
	"sync"

	"github.com/DusanKasan/cesium"
)

type canceller struct {
	cancelled bool
	onCancel  func()
	mux       sync.Mutex
}

func (c *canceller) Cancel() {
	c.mux.Lock()
	c.cancelled = true
	if c.onCancel != nil {
		c.onCancel()
	}
	c.mux.Unlock()
}

func (c *canceller) IsCancelled() bool {
	c.mux.Lock()
	x := c.cancelled
	c.mux.Unlock()

	return x
}

func (c *canceller) OnCancel(f func()) {
	c.mux.Lock()
	c.onCancel = f
	c.mux.Unlock()
}

type cancellable struct {
	cancel func()
}

func (c *cancellable) Cancel() {
	c.cancel()
}

type internalScheduler struct {
	schedule func(action func(cesium.Canceller)) cesium.Cancellable
}

func (is *internalScheduler) Schedule(action func(cesium.Canceller)) cesium.Cancellable {
	return is.schedule(action)
}

func SeparateGoroutineScheduler() cesium.Scheduler {
	c := make(chan func(), 10)

	go func() {
		for a := range c {
			a()
		}
	}()

	return &internalScheduler{
		schedule: func(action func(cesium.Canceller)) cesium.Cancellable {
			cc := &canceller{}

			c <- func() { action(cc) }

			return &cancellable{
				func() {
					cc.Cancel()
				},
			}
		},
	}
}
