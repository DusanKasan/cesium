package internal

import "github.com/DusanKasan/cesium"

func CountOperator(pub cesium.Publisher) cesium.Mono {
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
		onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
			p := CountProcessor()

			p.Subscribe(subscriber)
			return publisher.Subscribe(p)
		}

		return &Mono{onPublish}
	}
}
