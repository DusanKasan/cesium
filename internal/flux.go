package internal

import (
	"log"

	"github.com/DusanKasan/cesium"
)

type Flux struct {
	OnSubscribe func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription
}

func (f *Flux) Subscribe(subscriber cesium.Subscriber) cesium.Subscription {
	return f.OnSubscribe(subscriber, nil)
}

type ConditionalSubscriber interface {
	cesium.Subscriber
	OnNextIf(cesium.T) bool
}

func (f *Flux) Filter(filter func(t cesium.T) bool) cesium.Flux {
	return FluxFilterOperator(f, filter)
}

func (f *Flux) Map(mapper func(t cesium.T) cesium.T) cesium.Flux {
	return FluxMapOperator(f, mapper)
}

func (f *Flux) DoFinally(fn func()) cesium.Flux {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		p := DoFinallyProcessor(fn)

		subscription1 := p.Subscribe(subscriber)
		subscription2 := f.OnSubscribe(p, scheduler)
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

func (f *Flux) Count() cesium.Mono {
	return FluxCountOperator(f)
}

func (f *Flux) Reduce(fn func(cesium.T, cesium.T) cesium.T) cesium.Mono {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		p := ReduceProcessor(fn)

		subscription1 := p.Subscribe(subscriber)
		subscription2 := f.OnSubscribe(p, scheduler)
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

func (f *Flux) Scan(fn func(cesium.T, cesium.T) cesium.T) cesium.Flux {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		p := ScanProcessor(fn)

		subscription1 := p.Subscribe(subscriber)
		subscription2 := f.OnSubscribe(p, scheduler)
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

func (f *Flux) All(fn func(cesium.T) bool) cesium.Mono {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		p := AllProcessor(fn)

		subscription1 := p.Subscribe(subscriber)
		subscription2 := f.OnSubscribe(p, scheduler)
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

func (f *Flux) Any(fn func(cesium.T) bool) cesium.Mono {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		p := AnyProcessor(fn)

		subscription1 := p.Subscribe(subscriber)
		subscription2 := f.OnSubscribe(p, scheduler)

		subscription := &Subscription{
			CancelFunc: func() {
				subscription1.Cancel()
				subscription2.Cancel()
			},
			RequestFunc: func(n int64) {
				subscription1.Request(n)
			},
		}

		subscriber.OnSubscribe(subscription)
		return subscription

	}

	return &Mono{onPublish}
}

func (f *Flux) HasElements() cesium.Mono {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		p := HasElementsProcessor()

		subscription1 := p.Subscribe(subscriber)
		subscription2 := f.OnSubscribe(p, scheduler)
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

func (f *Flux) HasElement(element cesium.T) cesium.Mono {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		p := AnyProcessor(func(t cesium.T) bool {
			return t == element
		})

		subscription1 := p.Subscribe(subscriber)
		subscription2 := f.OnSubscribe(p, scheduler)
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

func (f *Flux) Log(logger *log.Logger) cesium.Flux {
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
		subscription2 := f.OnSubscribe(p, scheduler)
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

func (f *Flux) DoOnNext(fn func(cesium.T)) cesium.Flux {
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
		subscription2 := f.OnSubscribe(p, scheduler)
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

func (f *Flux) DoOnError(fn func(error)) cesium.Flux {
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
		subscription2 := f.OnSubscribe(p, scheduler)
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

func (f *Flux) DoOnCancel(fn func()) cesium.Flux {
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
		subscription2 := f.OnSubscribe(p, scheduler)
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

func (f *Flux) DoOnSubscribe(fn func(cesium.Subscription)) cesium.Flux {
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
		subscription2 := f.OnSubscribe(p, scheduler)
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

func (f *Flux) DoOnRequest(fn func(int64)) cesium.Flux {
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
		subscription2 := f.OnSubscribe(p, scheduler)
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

func (f *Flux) DoOnTerminate(fn func()) cesium.Flux {
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
		subscription2 := f.OnSubscribe(p, scheduler)
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

func (f *Flux) DoOnComplete(fn func()) cesium.Flux {
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
		subscription2 := f.OnSubscribe(p, scheduler)
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

func (f *Flux) DoAfterTerminate(fn func()) cesium.Flux {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		p := DoAfterTerminateProcessor(fn)

		subscription1 := p.Subscribe(subscriber)
		subscription2 := f.OnSubscribe(p, scheduler)
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

func (f *Flux) Handle(fn func(cesium.T, cesium.SynchronousSink)) cesium.Flux {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		p := HandleProcessor(fn)

		subscription1 := p.Subscribe(subscriber)
		subscription2 := f.OnSubscribe(p, scheduler)
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

func (f *Flux) Concat(publishers cesium.Publisher /*<cesium.Publisher>*/) cesium.Flux {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		p := ConcatProcessor(publishers)

		subscription1 := p.Subscribe(subscriber)
		subscription2 := f.OnSubscribe(p, scheduler)
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

func (f *Flux) ConcatWith(publishers ...cesium.Publisher) cesium.Flux {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		var pubs []cesium.T
		for _, pub := range publishers {
			pubs = append(pubs, pub.(cesium.T))
		}

		p := ConcatProcessor(FluxFromSlice(pubs))

		subscription1 := p.Subscribe(subscriber)
		subscription2 := f.OnSubscribe(p, scheduler)
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

func (f *Flux) DistinctUntilChanged() cesium.Flux {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		p := DistinctUntilChangedProcessor()

		subscription1 := p.Subscribe(subscriber)
		subscription2 := f.OnSubscribe(p, scheduler)
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

func (f *Flux) Take(n int64) cesium.Flux {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		p := TakeProcessor(n)

		subscription1 := p.Subscribe(subscriber)
		subscription2 := f.OnSubscribe(p, scheduler)
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
