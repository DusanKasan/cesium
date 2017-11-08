package internal

import (
	"fmt"

	"github.com/DusanKasan/cesium"
)

func FluxCountOperator(pub cesium.Publisher) cesium.Mono {
	switch publisher := pub.(type) {
	case ScalarCallable:
		_, ok := publisher.Get()
		if ok {
			return MonoJust(int64(1))
		}

		return MonoJust(int64(0))
	case *Flux:
		onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
			p := CountProcessor()

			subscription1 := p.Subscribe(subscriber)
			subscription2 := publisher.OnSubscribe(p, scheduler)
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
	default:
		panic(fmt.Sprintf("invalid publisher type: %v", pub))
		//	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		//		requested := int64(0)
		//		mux := sync.Mutex{}
		//		var sub cesium.Subscription
		//
		//		p := CountProcessor()
		//
		//		cancellable := scheduler.Schedule(func(c cesium.Canceller) {
		//			if c.IsCancelled() {
		//				return
		//			}
		//
		//			p.Subscribe(subscriber)
		//			s := publisher.Subscribe(p)
		//			mux.Lock()
		//			sub = s
		//			if requested > 0 {
		//				sub.Request(requested)
		//			}
		//			mux.Unlock()
		//		})
		//
		//		return &Subscription{
		//			func() {
		//				sub.Cancel()
		//				cancellable.Cancel()
		//			},
		//			func(i int64) {
		//				mux.Lock()
		//				if sub == nil {
		//					requested = requested + i
		//				} else {
		//					sub.Request(i)
		//				}
		//				mux.Unlock()
		//			},
		//		}
		//	}
		//
		//	return &Mono{onPublish}
	}
}

func FluxFilterOperator(pub cesium.Publisher, f func(cesium.T) bool) cesium.Flux {
	switch publisher := pub.(type) {
	case *ScalarFlux:
		t, ok := publisher.Get()
		if ok && f(t) {
			return publisher
		}

		return fluxFromCallable(func() (cesium.T, bool) {
			return nil, false
		})
	case *Flux:
		onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
			p := FilterProcessor(f)

			subscription1 := p.Subscribe(subscriber)
			subscription2 := publisher.OnSubscribe(p, scheduler)
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
	default:
		panic(fmt.Sprintf("invalid publisher type: %v", pub))
	}
}

func FluxMapOperator(pub cesium.Publisher, f func(t cesium.T) cesium.T) cesium.Flux {
	switch publisher := pub.(type) {
	case *ScalarFlux:
		t, ok := publisher.Get()

		if ok {
			t = f(t)
		}

		return fluxFromCallable(func() (cesium.T, bool) {
			return t, ok
		})
	case *Flux:
		onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
			p := MapProcessor(f)

			subscription1 := p.Subscribe(subscriber)
			subscription2 := publisher.OnSubscribe(p, scheduler)
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
	default:
		panic(fmt.Sprintf("invalid publisher type: %v", pub))
	}
}

func MonoFilterOperator(pub cesium.Publisher, f func(t cesium.T) bool) cesium.Mono {
	switch publisher := pub.(type) {
	case *ScalarMono:
		t, ok := publisher.Get()
		if ok && f(t) {
			return publisher
		}

		return monoFromCallable(func() (cesium.T, bool) {
			return nil, false
		})
	case *Mono:
		onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
			p := FilterProcessor(f)

			subscription1 := p.Subscribe(subscriber)
			subscription2 := publisher.OnSubscribe(p, scheduler)
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
	default:
		panic(fmt.Sprintf("invalid publisher type: %v", pub))
	}
}

func MonoMapOperator(pub cesium.Publisher, f func(t cesium.T) cesium.T) cesium.Mono {
	switch publisher := pub.(type) {
	case *ScalarMono:
		t, ok := publisher.Get()

		if ok {
			t = f(t)
		}

		return monoFromCallable(func() (cesium.T, bool) {
			return t, ok
		})
	case *Mono:
		onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
			p := MapProcessor(f)

			subscription1 := p.Subscribe(subscriber)
			subscription2 := publisher.OnSubscribe(p, scheduler)
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
	default:
		panic(fmt.Sprintf("invalid publisher type: %v", pub))
	}
}
