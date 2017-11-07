package internal

import (
	"log"

	"github.com/DusanKasan/cesium"
)

type ScalarMono struct {
	cesium.Mono
	get func() (cesium.T, bool)
}

func (s *ScalarMono) Get() (cesium.T, bool) {
	return s.get()
}

type Mono struct {
	OnSubscribe func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription
}

func (m *Mono) Subscribe(subscriber cesium.Subscriber) cesium.Subscription {
	return m.OnSubscribe(subscriber, nil)
}

func (m *Mono) Filter(filter func(t cesium.T) bool) cesium.Mono {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		p := FilterProcessor(filter)

		subscription1 := p.Subscribe(subscriber)
		subscription2 := m.OnSubscribe(p, scheduler)
		p.OnSubscribe(subscription2)

		sub := &Subscription{
			CancelFunc: func() {
				subscription1.Cancel()
				subscription2.Cancel()
			},
			RequestFunc: func(n int64) {
				subscription1.Request(n)
			},
		}

		subscriber.OnSubscribe(sub)
		return sub
	}

	return &Mono{onPublish}
}

func (m *Mono) Map(mapper func(t cesium.T) cesium.T) cesium.Mono {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		p := MapProcessor(mapper)

		subscription1 := p.Subscribe(subscriber)
		subscription2 := m.OnSubscribe(p, scheduler)
		p.OnSubscribe(subscription2)

		sub := &Subscription{
			CancelFunc: func() {
				subscription1.Cancel()
				subscription2.Cancel()
			},
			RequestFunc: func(n int64) {
				subscription1.Request(n)
			},
		}

		subscriber.OnSubscribe(sub)
		return sub
	}

	return &Mono{onPublish}
}

func (m *Mono) DoFinally(fn func()) cesium.Mono {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		p := DoFinallyProcessor(fn)

		subscription1 := p.Subscribe(subscriber)
		subscription2 := m.OnSubscribe(p, scheduler)
		p.OnSubscribe(subscription2)

		sub := &Subscription{
			CancelFunc: func() {
				subscription1.Cancel()
				subscription2.Cancel()
			},
			RequestFunc: func(n int64) {
				subscription1.Request(n)
			},
		}

		subscriber.OnSubscribe(sub)
		return sub
	}

	return &Mono{onPublish}
}

func (m *Mono) Log(logger *log.Logger) cesium.Mono {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		p := DoProcessor(
			func(t cesium.T) {
				logger.Printf("Next: %#v\n", t)
			},
			func() {
				logger.Println("Complete")
			},
			func(err error) {
				logger.Printf("Error: %#v\n", err)
			},
			func() {
				logger.Printf("Cancel")
			},
			func(n int64) {
				logger.Printf("Request: %#v\n", n)
			},
			func(cesium.Subscription) {},
		)

		logger.Printf("Subscribed")
		subscription1 := p.Subscribe(subscriber)
		subscription2 := m.OnSubscribe(p, scheduler)

		sub := &Subscription{
			CancelFunc: func() {
				subscription1.Cancel()
				subscription2.Cancel()
			},
			RequestFunc: func(n int64) {
				subscription1.Request(n)
			},
		}

		subscriber.OnSubscribe(sub)
		return sub
	}

	return &Mono{onPublish}
}

func (m *Mono) DoOnNext(fn func(cesium.T)) cesium.Mono {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		p := DoProcessor(
			fn,
			func() {},
			func(error) {},
			func() {},
			func(int64) {},
			func(cesium.Subscription) {},
		)

		subscription1 := p.Subscribe(subscriber)
		subscription2 := m.OnSubscribe(p, scheduler)
		p.OnSubscribe(subscription2)

		sub := &Subscription{
			CancelFunc: func() {
				subscription1.Cancel()
				subscription2.Cancel()
			},
			RequestFunc: func(n int64) {
				subscription1.Request(n)
			},
		}

		subscriber.OnSubscribe(sub)
		return sub
	}

	return &Mono{onPublish}
}

func (m *Mono) DoOnError(fn func(error)) cesium.Mono {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		p := DoProcessor(
			func(cesium.T) {},
			func() {},
			func(err error) {
				fn(err)
			},
			func() {},
			func(int64) {},
			func(cesium.Subscription) {},
		)

		subscription1 := p.Subscribe(subscriber)
		subscription2 := m.OnSubscribe(p, scheduler)
		p.OnSubscribe(subscription2)

		sub := &Subscription{
			CancelFunc: func() {
				subscription1.Cancel()
				subscription2.Cancel()
			},
			RequestFunc: func(n int64) {
				subscription1.Request(n)
			},
		}

		subscriber.OnSubscribe(sub)
		return sub
	}

	return &Mono{onPublish}
}

func (m *Mono) DoOnCancel(fn func()) cesium.Mono {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		p := DoProcessor(
			func(cesium.T) {},
			func() {},
			func(err error) {},
			func() {
				fn()
			},
			func(int64) {},
			func(cesium.Subscription) {},
		)

		subscription1 := p.Subscribe(subscriber)
		subscription2 := m.OnSubscribe(p, scheduler)
		p.OnSubscribe(subscription2)

		sub := &Subscription{
			CancelFunc: func() {
				subscription1.Cancel()
				subscription2.Cancel()
			},
			RequestFunc: func(n int64) {
				subscription1.Request(n)
			},
		}

		subscriber.OnSubscribe(sub)
		return sub
	}

	return &Mono{onPublish}
}

func (m *Mono) DoOnSubscribe(fn func(cesium.Subscription)) cesium.Mono {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		p := DoProcessor(
			func(cesium.T) {},
			func() {},
			func(err error) {},
			func() {},
			func(int64) {},
			func(s cesium.Subscription) {
				fn(s)
			},
		)

		subscription1 := p.Subscribe(subscriber)
		subscription2 := m.OnSubscribe(p, scheduler)
		p.OnSubscribe(subscription2)

		sub := &Subscription{
			CancelFunc: func() {
				subscription1.Cancel()
				subscription2.Cancel()
			},
			RequestFunc: func(n int64) {
				subscription1.Request(n)
			},
		}

		subscriber.OnSubscribe(sub)
		return sub
	}

	return &Mono{onPublish}
}

func (m *Mono) DoOnRequest(fn func(int64)) cesium.Mono {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		p := DoProcessor(
			func(cesium.T) {},
			func() {},
			func(err error) {},
			func() {},
			func(n int64) {
				fn(n)
			},
			func(s cesium.Subscription) {},
		)

		subscription1 := p.Subscribe(subscriber)
		subscription2 := m.OnSubscribe(p, scheduler)
		p.OnSubscribe(subscription2)

		sub := &Subscription{
			CancelFunc: func() {
				subscription1.Cancel()
				subscription2.Cancel()
			},
			RequestFunc: func(n int64) {
				subscription1.Request(n)
			},
		}

		subscriber.OnSubscribe(sub)
		return sub
	}

	return &Mono{onPublish}
}

func (m *Mono) DoOnTerminate(fn func()) cesium.Mono {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		p := DoProcessor(
			func(cesium.T) {},
			func() {
				fn()
			},
			func(err error) {
				fn()
			},
			func() {},
			func(int64) {},
			func(cesium.Subscription) {},
		)

		subscription1 := p.Subscribe(subscriber)
		subscription2 := m.OnSubscribe(p, scheduler)
		p.OnSubscribe(subscription2)

		sub := &Subscription{
			CancelFunc: func() {
				subscription1.Cancel()
				subscription2.Cancel()
			},
			RequestFunc: func(n int64) {
				subscription1.Request(n)
			},
		}

		subscriber.OnSubscribe(sub)
		return sub
	}

	return &Mono{onPublish}
}

func (m *Mono) DoOnSuccess(fn func()) cesium.Mono {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		p := DoProcessor(
			func(cesium.T) {},
			func() {
				fn()
			},
			func(err error) {},
			func() {},
			func(int64) {},
			func(cesium.Subscription) {},
		)

		subscription1 := p.Subscribe(subscriber)
		subscription2 := m.OnSubscribe(p, scheduler)
		p.OnSubscribe(subscription2)

		sub := &Subscription{
			CancelFunc: func() {
				subscription1.Cancel()
				subscription2.Cancel()
			},
			RequestFunc: func(n int64) {
				subscription1.Request(n)
			},
		}

		subscriber.OnSubscribe(sub)
		return sub
	}

	return &Mono{onPublish}
}

func (m *Mono) DoAfterTerminate(fn func()) cesium.Mono {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		p := DoAfterTerminateProcessor(fn)

		subscription1 := p.Subscribe(subscriber)
		subscription2 := m.OnSubscribe(p, scheduler)
		p.OnSubscribe(subscription2)

		sub := &Subscription{
			CancelFunc: func() {
				subscription1.Cancel()
				subscription2.Cancel()
			},
			RequestFunc: func(n int64) {
				subscription1.Request(n)
			},
		}

		subscriber.OnSubscribe(sub)
		return sub
	}

	return &Mono{onPublish}
}

func (m *Mono) ConcatWith(publishers ...cesium.Publisher) cesium.Flux {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		var pubs []cesium.T
		for _, pub := range publishers {
			pubs = append(pubs, pub.(cesium.T))
		}

		p := ConcatProcessor(FluxFromSlice(pubs))

		subscription1 := p.Subscribe(subscriber)
		subscription2 := m.OnSubscribe(p, scheduler)
		p.OnSubscribe(subscription2)

		sub := &Subscription{
			CancelFunc: func() {
				subscription1.Cancel()
				subscription2.Cancel()
			},
			RequestFunc: func(n int64) {
				subscription1.Request(n)
			},
		}

		subscriber.OnSubscribe(sub)
		return sub
	}

	return &Flux{onPublish}
}
