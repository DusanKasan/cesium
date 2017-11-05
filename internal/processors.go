package internal

import (
	"sync"

	"math"

	"github.com/DusanKasan/cesium"
)

//NOTE: processor can only subscribe to one publisher

type processor struct {
	onNext      func(cesium.T)
	onComplete  func()
	onError     func(error)
	onSubscribe func(cesium.Subscription)
	subscribe   func(cesium.Subscriber) cesium.Subscription
}

func (p *processor) OnNext(t cesium.T) {
	p.onNext(t)
}

func (p *processor) OnComplete() {
	p.onComplete()
}

func (p *processor) OnError(err error) {
	p.onError(err)
}

func (p *processor) OnSubscribe(s cesium.Subscription) {
	p.onSubscribe(s)
}

func (p *processor) Subscribe(s cesium.Subscriber) cesium.Subscription {
	return p.subscribe(s)
}

func FilterProcessor(f func(cesium.T) bool) cesium.Processor {
	var subscriber cesium.Subscriber
	var subscription cesium.Subscription
	subscriberMux := sync.Mutex{}
	subscribtionMux := sync.Mutex{}

	return &processor{
		subscribe: func(s cesium.Subscriber) cesium.Subscription {
			sub := &Subscription{
				CancelFunc: func() {
					subscribtionMux.Lock()
					if subscription != nil {
						subscription.Cancel()
					}
					subscribtionMux.Unlock()
				},
				RequestFunc: func(n int64) {
					subscribtionMux.Lock()
					if subscription != nil {
						subscription.Request(n)
					}
					subscribtionMux.Unlock()
				},
			}

			subscriberMux.Lock()
			subscriber = s
			subscriber.OnSubscribe(subscription)
			subscriberMux.Unlock()

			return sub
		},
		onSubscribe: func(s cesium.Subscription) {
			subscribtionMux.Lock()
			subscription = s
			subscribtionMux.Unlock()
		},
		onNext: func(t cesium.T) {
			if f(t) {
				subscriberMux.Lock()
				subscriber.OnNext(t)
				subscriberMux.Unlock()
			} else {
				subscribtionMux.Lock()
				subscription.Request(1)
				subscribtionMux.Unlock()
			}
		},
		onComplete: func() {
			subscriberMux.Lock()
			subscriber.OnComplete()
			subscriberMux.Unlock()
		},
		onError: func(err error) {
			subscriberMux.Lock()
			subscriber.OnError(err)
			subscriberMux.Unlock()
		},
	}
}

func MapProcessor(f func(cesium.T) cesium.T) cesium.Processor {
	var subscriber cesium.Subscriber
	var subscription cesium.Subscription
	subscriberMux := sync.Mutex{}
	subscribtionMux := sync.Mutex{}

	return &processor{
		subscribe: func(s cesium.Subscriber) cesium.Subscription {
			sub := &Subscription{
				CancelFunc: func() {
					subscribtionMux.Lock()
					if subscription != nil {
						subscription.Cancel()
					}
					subscribtionMux.Unlock()
				},
				RequestFunc: func(n int64) {
					subscribtionMux.Lock()
					if subscription != nil {
						subscription.Request(n)
					}
					subscribtionMux.Unlock()
				},
			}

			subscriberMux.Lock()
			subscriber = s
			subscriber.OnSubscribe(subscription)
			subscriberMux.Unlock()

			return sub
		},
		onSubscribe: func(s cesium.Subscription) {
			subscribtionMux.Lock()
			subscription = s
			subscribtionMux.Unlock()
		},
		onNext: func(t cesium.T) {
			subscriberMux.Lock()
			subscriber.OnNext(f(t))
			subscriberMux.Unlock()
		},
		onComplete: func() {
			subscriberMux.Lock()
			subscriber.OnComplete()
			subscriberMux.Unlock()
		},
		onError: func(err error) {
			subscriberMux.Lock()
			subscriber.OnError(err)
			subscriberMux.Unlock()
		},
	}
}

func DoFinallyProcessor(f func()) cesium.Processor {
	var subscriber cesium.Subscriber
	var subscription cesium.Subscription
	subscriberMux := sync.Mutex{}
	subscribtionMux := sync.Mutex{}

	return &processor{
		subscribe: func(s cesium.Subscriber) cesium.Subscription {
			sub := &Subscription{
				CancelFunc: func() {
					subscribtionMux.Lock()
					if subscription != nil {
						subscription.Cancel()
						f()
					}
					subscribtionMux.Unlock()
				},
				RequestFunc: func(n int64) {
					subscribtionMux.Lock()
					if subscription != nil {
						subscription.Request(n)
					}
					subscribtionMux.Unlock()
				},
			}

			subscriberMux.Lock()
			subscriber = s
			subscriber.OnSubscribe(subscription)
			subscriberMux.Unlock()

			return sub
		},
		onSubscribe: func(s cesium.Subscription) {
			subscribtionMux.Lock()
			subscription = s
			subscribtionMux.Unlock()
		},
		onNext: func(t cesium.T) {
			subscriberMux.Lock()
			subscriber.OnNext(t)
			subscriberMux.Unlock()
		},
		onComplete: func() {
			subscriberMux.Lock()
			subscriber.OnComplete()
			subscriberMux.Unlock()
			f()
		},
		onError: func(err error) {
			subscriberMux.Lock()
			subscriber.OnError(err)
			subscriberMux.Unlock()
			f()
		},
	}
}

func CountProcessor() cesium.Processor {
	var subscriber cesium.Subscriber
	var subscription cesium.Subscription
	subscriberMux := sync.Mutex{}
	subscribtionMux := sync.Mutex{}

	count := int64(0)

	return &processor{
		subscribe: func(s cesium.Subscriber) cesium.Subscription {
			sub := &Subscription{
				CancelFunc: func() {
					subscribtionMux.Lock()
					if subscription != nil {
						subscription.Cancel()
					}
					subscribtionMux.Unlock()
				},
				RequestFunc: func(n int64) {
					subscribtionMux.Lock()
					if subscription != nil {
						subscription.Request(1)
					}
					subscribtionMux.Unlock()
				},
			}

			subscriberMux.Lock()
			subscriber = s
			//subscriber.OnSubscribe(subscription)
			subscriberMux.Unlock()

			return sub
		},
		onSubscribe: func(s cesium.Subscription) {
			subscribtionMux.Lock()
			subscription = s
			subscribtionMux.Unlock()
		},
		onNext: func(t cesium.T) {
			count++
			subscription.Request(1)
		},
		onComplete: func() {
			subscriberMux.Lock()
			subscriber.OnNext(count)
			subscriber.OnComplete()
			subscriberMux.Unlock()
		},
		onError: func(err error) {
			subscriberMux.Lock()
			subscriber.OnError(err)
			subscriberMux.Unlock()
		},
	}
}

func ReduceProcessor(f func(cesium.T, cesium.T) cesium.T) cesium.Processor {
	var subscriber cesium.Subscriber
	var subscription cesium.Subscription
	subscriberMux := sync.Mutex{}
	subscribtionMux := sync.Mutex{}

	firstItem := true
	var previousItem cesium.T
	mux := sync.Mutex{}

	return &processor{
		subscribe: func(s cesium.Subscriber) cesium.Subscription {
			sub := &Subscription{
				CancelFunc: func() {
					subscribtionMux.Lock()
					if subscription != nil {
						subscription.Cancel()
					}
					subscribtionMux.Unlock()
				},
				RequestFunc: func(n int64) {
					subscribtionMux.Lock()
					if subscription != nil {
						subscription.Request(1)
					}
					subscribtionMux.Unlock()
				},
			}

			subscriberMux.Lock()
			subscriber = s
			subscriber.OnSubscribe(subscription)
			subscriberMux.Unlock()

			return sub
		},
		onSubscribe: func(s cesium.Subscription) {
			subscribtionMux.Lock()
			subscription = s
			subscribtionMux.Unlock()
		},
		onNext: func(t cesium.T) {
			mux.Lock()
			if firstItem {
				firstItem = false
				previousItem = t
			} else {
				previousItem = f(previousItem, t)
			}
			mux.Unlock()
			subscription.Request(1)
		},
		onComplete: func() {
			subscriberMux.Lock()
			subscriber.OnNext(previousItem)
			subscriber.OnComplete()
			subscriberMux.Unlock()
		},
		onError: func(err error) {
			subscriberMux.Lock()
			subscriber.OnError(err)
			subscriberMux.Unlock()
		},
	}
}

func ScanProcessor(f func(cesium.T, cesium.T) cesium.T) cesium.Processor {
	var subscriber cesium.Subscriber
	var subscription cesium.Subscription
	subscriberMux := sync.Mutex{}
	subscribtionMux := sync.Mutex{}

	firstItem := true
	var previousItem cesium.T
	mux := sync.Mutex{}

	return &processor{
		subscribe: func(s cesium.Subscriber) cesium.Subscription {
			sub := &Subscription{
				CancelFunc: func() {
					subscribtionMux.Lock()
					if subscription != nil {
						subscription.Cancel()
					}
					subscribtionMux.Unlock()
				},
				RequestFunc: func(n int64) {
					subscribtionMux.Lock()
					if subscription != nil {
						subscription.Request(1)
					}
					subscribtionMux.Unlock()
				},
			}

			subscriberMux.Lock()
			subscriber = s
			subscriber.OnSubscribe(subscription)
			subscriberMux.Unlock()

			return sub
		},
		onSubscribe: func(s cesium.Subscription) {
			subscribtionMux.Lock()
			subscription = s
			subscribtionMux.Unlock()
		},
		onNext: func(t cesium.T) {
			mux.Lock()
			if firstItem {
				firstItem = false
				previousItem = t
			} else {
				previousItem = f(previousItem, t)
			}
			p := previousItem
			mux.Unlock()
			subscriber.OnNext(p)
			subscription.Request(1)
		},
		onComplete: func() {
			subscriberMux.Lock()
			subscriber.OnComplete()
			subscriberMux.Unlock()
		},
		onError: func(err error) {
			subscriberMux.Lock()
			subscriber.OnError(err)
			subscriberMux.Unlock()
		},
	}
}

func AllProcessor(f func(cesium.T) bool) cesium.Processor {
	var subscriber cesium.Subscriber
	var subscription cesium.Subscription
	subscriberMux := sync.Mutex{}
	subscribtionMux := sync.Mutex{}

	closed := false

	return &processor{
		subscribe: func(s cesium.Subscriber) cesium.Subscription {
			sub := &Subscription{
				CancelFunc: func() {
					subscribtionMux.Lock()
					if subscription != nil {
						subscription.Cancel()
					}
					subscribtionMux.Unlock()
				},
				RequestFunc: func(n int64) {
					subscribtionMux.Lock()
					if subscription != nil {
						subscription.Request(1)
					}
					subscribtionMux.Unlock()
				},
			}

			subscriberMux.Lock()
			subscriber = s
			subscriber.OnSubscribe(subscription)
			subscriberMux.Unlock()

			return sub
		},
		onSubscribe: func(s cesium.Subscription) {
			subscribtionMux.Lock()
			subscription = s
			subscribtionMux.Unlock()
		},
		onNext: func(t cesium.T) {
			if f(t) {
				subscribtionMux.Lock()
				subscription.Request(1)
				subscribtionMux.Unlock()
				return
			} else {
				subscriberMux.Lock()
				closed = true
				subscriber.OnNext(false)
				subscriber.OnComplete()
				subscriberMux.Unlock()

				subscribtionMux.Lock()
				subscription.Cancel()
				subscribtionMux.Unlock()
			}
		},
		onComplete: func() {
			subscriberMux.Lock()
			if !closed {
				subscriber.OnNext(true)
				subscriber.OnComplete()
			}
			subscriberMux.Unlock()
		},
		onError: func(err error) {
			subscriberMux.Lock()
			if !closed {
				subscriber.OnError(err)
			}
			subscriberMux.Unlock()
		},
	}
}

func AnyProcessor(f func(cesium.T) bool) cesium.Processor {
	var subscriber cesium.Subscriber
	var subscription cesium.Subscription
	subscriberMux := sync.Mutex{}
	subscribtionMux := sync.Mutex{}

	closed := false

	return &processor{
		subscribe: func(s cesium.Subscriber) cesium.Subscription {
			sub := &Subscription{
				CancelFunc: func() {
					subscribtionMux.Lock()
					if subscription != nil {
						subscription.Cancel()
					}
					subscribtionMux.Unlock()
				},
				RequestFunc: func(n int64) {
					subscribtionMux.Lock()
					if subscription != nil {
						subscription.Request(1)
					}
					subscribtionMux.Unlock()
				},
			}

			subscriberMux.Lock()
			subscriber = s
			subscriber.OnSubscribe(subscription)
			subscriberMux.Unlock()

			return sub
		},
		onSubscribe: func(s cesium.Subscription) {
			subscribtionMux.Lock()
			subscription = s
			subscribtionMux.Unlock()
		},
		onNext: func(t cesium.T) {
			if !f(t) {
				subscribtionMux.Lock()
				subscription.Request(1)
				subscribtionMux.Unlock()
				return
			} else {
				subscriberMux.Lock()
				closed = true
				subscriber.OnNext(true)
				subscriber.OnComplete()
				subscriberMux.Unlock()

				subscribtionMux.Lock()
				subscription.Cancel()
				subscribtionMux.Unlock()
			}
		},
		onComplete: func() {
			subscriberMux.Lock()
			if !closed {
				subscriber.OnNext(false)
				subscriber.OnComplete()
			}
			subscriberMux.Unlock()
		},
		onError: func(err error) {
			subscriberMux.Lock()
			if !closed {
				subscriber.OnError(err)
			}
			subscriberMux.Unlock()
		},
	}
}

func HasElementsProcessor() cesium.Processor {
	var subscriber cesium.Subscriber
	var subscription cesium.Subscription
	subscriberMux := sync.Mutex{}
	subscribtionMux := sync.Mutex{}

	closed := false

	return &processor{
		subscribe: func(s cesium.Subscriber) cesium.Subscription {
			sub := &Subscription{
				CancelFunc: func() {
					subscribtionMux.Lock()
					if subscription != nil {
						subscription.Cancel()
					}
					subscribtionMux.Unlock()
				},
				RequestFunc: func(n int64) {
					subscribtionMux.Lock()
					if subscription != nil {
						subscription.Request(1)
					}
					subscribtionMux.Unlock()
				},
			}

			subscriberMux.Lock()
			subscriber = s
			subscriber.OnSubscribe(subscription)
			subscriberMux.Unlock()

			return sub
		},
		onSubscribe: func(s cesium.Subscription) {
			subscribtionMux.Lock()
			subscription = s
			subscribtionMux.Unlock()
		},
		onNext: func(t cesium.T) {
			subscriberMux.Lock()
			closed = true
			subscriber.OnNext(true)
			subscriber.OnComplete()
			subscriberMux.Unlock()

			subscribtionMux.Lock()
			subscription.Cancel()
			subscribtionMux.Unlock()
		},
		onComplete: func() {
			subscriberMux.Lock()
			if !closed {
				subscriber.OnNext(false)
				subscriber.OnComplete()
			}
			subscriberMux.Unlock()
		},
		onError: func(err error) {
			subscriberMux.Lock()
			if !closed {
				subscriber.OnError(err)
			}
			subscriberMux.Unlock()
		},
	}
}

func DoProcessor(
	onNext func(cesium.T),
	onComplete func(),
	onError func(error),
	onCancel func(),
	onRequest func(int64),
	onSubscribe func(cesium.Subscription),
) cesium.Processor {
	var subscriber cesium.Subscriber
	var subscription cesium.Subscription
	subscriberMux := sync.Mutex{}
	subscribtionMux := sync.Mutex{}

	return &processor{
		subscribe: func(s cesium.Subscriber) cesium.Subscription {
			sub := &Subscription{
				CancelFunc: func() {
					subscribtionMux.Lock()
					onCancel()
					if subscription != nil {
						subscription.Cancel()
					}
					subscribtionMux.Unlock()
				},
				RequestFunc: func(n int64) {
					subscribtionMux.Lock()
					onRequest(n)
					if subscription != nil {
						subscription.Request(n)
					}
					subscribtionMux.Unlock()
				},
			}

			subscriberMux.Lock()
			subscriber = s
			subscriber.OnSubscribe(subscription)
			subscriberMux.Unlock()

			return sub
		},
		onSubscribe: func(s cesium.Subscription) {
			subscribtionMux.Lock()
			onSubscribe(s)
			subscription = s
			subscribtionMux.Unlock()
		},
		onNext: func(t cesium.T) {
			subscriberMux.Lock()
			onNext(t)
			subscriber.OnNext(t)
			subscriberMux.Unlock()
		},
		onComplete: func() {
			subscriberMux.Lock()
			onComplete()
			subscriber.OnComplete()
			subscriberMux.Unlock()
		},
		onError: func(err error) {
			subscriberMux.Lock()
			onError(err)
			subscriber.OnError(err)
			subscriberMux.Unlock()
		},
	}
}

func DoAfterTerminateProcessor(fn func()) cesium.Processor {
	var subscriber cesium.Subscriber
	var subscription cesium.Subscription
	subscriberMux := sync.Mutex{}
	subscribtionMux := sync.Mutex{}

	return &processor{
		subscribe: func(s cesium.Subscriber) cesium.Subscription {
			sub := &Subscription{
				CancelFunc: func() {
					subscribtionMux.Lock()
					if subscription != nil {
						subscription.Cancel()
					}
					subscribtionMux.Unlock()
				},
				RequestFunc: func(n int64) {
					subscribtionMux.Lock()
					if subscription != nil {
						subscription.Request(n)
					}
					subscribtionMux.Unlock()
				},
			}

			subscriberMux.Lock()
			subscriber = s
			subscriber.OnSubscribe(subscription)
			subscriberMux.Unlock()

			return sub
		},
		onSubscribe: func(s cesium.Subscription) {
			subscribtionMux.Lock()
			subscription = s
			subscribtionMux.Unlock()
		},
		onNext: func(t cesium.T) {
			subscriberMux.Lock()
			subscriber.OnNext(t)
			subscriberMux.Unlock()
		},
		onComplete: func() {
			subscriberMux.Lock()
			subscriber.OnComplete()
			subscriberMux.Unlock()
			fn()
		},
		onError: func(err error) {
			subscriberMux.Lock()
			subscriber.OnError(err)
			subscriberMux.Unlock()
			fn()
		},
	}
}

func HandleProcessor(fn func(cesium.T, cesium.SynchronousSink)) cesium.Processor {
	var subscriber cesium.Subscriber
	var subscription cesium.Subscription
	subscriberMux := sync.Mutex{}
	subscribtionMux := sync.Mutex{}

	terminated := false

	return &processor{
		subscribe: func(s cesium.Subscriber) cesium.Subscription {
			sub := &Subscription{
				CancelFunc: func() {
					subscribtionMux.Lock()
					if subscription != nil {
						subscription.Cancel()
					}
					subscribtionMux.Unlock()
				},
				RequestFunc: func(n int64) {
					subscribtionMux.Lock()
					if subscription != nil {
						subscription.Request(n)
					}
					subscribtionMux.Unlock()
				},
			}

			subscriberMux.Lock()
			subscriber = s
			subscriber.OnSubscribe(subscription)
			subscriberMux.Unlock()

			return sub
		},
		onSubscribe: func(s cesium.Subscription) {
			subscribtionMux.Lock()
			subscription = s
			subscribtionMux.Unlock()
		},
		onNext: func(t cesium.T) {
			subscriberMux.Lock()
			if terminated {
				subscriberMux.Unlock()
				return
			}

			sink := &SynchronousSink{}
			fn(t, sink)
			emission := sink.GetEmission()
			switch emission.EventType {
			case "next":
				subscriber.OnNext(emission.Value)
			case "complete":
				subscriber.OnComplete()
				terminated = true
				subscribtionMux.Lock()
				subscription.Cancel()
				subscribtionMux.Unlock()
			case "error":
				subscriber.OnError(emission.Err)
				terminated = true
				subscribtionMux.Lock()
				subscription.Cancel()
				subscribtionMux.Unlock()
			default:
				subscriber.OnError(cesium.NoEmissionOnSynchronousSinkError)
				terminated = true
				subscribtionMux.Lock()
				subscription.Cancel()
				subscribtionMux.Unlock()
			}
			subscriberMux.Unlock()
		},
		onComplete: func() {
			subscriberMux.Lock()
			subscriber.OnComplete()
			subscriberMux.Unlock()
		},
		onError: func(err error) {
			subscriberMux.Lock()
			subscriber.OnError(err)
			subscriberMux.Unlock()
		},
	}
}

func ConcatProcessor(publishers cesium.Publisher /*<cesium.Publisher>*/) cesium.Processor {
	var subscriber cesium.Subscriber
	var subscription cesium.Subscription
	subscriberMux := sync.Mutex{}
	subscribtionMux := sync.Mutex{}

	var proc *processor

	mux := sync.Mutex{}
	publishersComplete := false
	pendingRequests := int64(0)
	unbounded := false

	p := DoObserver(
		func(t cesium.T) {
			s := t.(cesium.Publisher).Subscribe(proc)

			if !unbounded {
				s.Request(pendingRequests)
			}
		},
		func() {
			mux.Lock()
			publishersComplete = true
			mux.Unlock()
		},
		func(err error) {
			MonoError(err).Subscribe(proc).Request(1)
		},
	)

	var publishersSubscription cesium.Subscription
	publishersSubscription = publishers.Subscribe(p)

	proc = &processor{
		subscribe: func(s cesium.Subscriber) cesium.Subscription {
			sub := &Subscription{
				CancelFunc: func() {
					subscribtionMux.Lock()
					if subscription != nil {
						subscription.Cancel()
					}
					subscribtionMux.Unlock()
				},
				RequestFunc: func(n int64) {
					subscribtionMux.Lock()
					if n == math.MaxInt64 {
						unbounded = true
					}

					if !unbounded {
						pendingRequests = pendingRequests + n
					}

					if subscription != nil {
						subscription.Request(n)
					}
					subscribtionMux.Unlock()
				},
			}

			subscriberMux.Lock()
			subscriber = s
			subscriber.OnSubscribe(subscription)
			subscriberMux.Unlock()

			return sub
		},
		onSubscribe: func(s cesium.Subscription) {
			subscribtionMux.Lock()
			subscription = s
			subscribtionMux.Unlock()
		},
		onNext: func(t cesium.T) {
			subscriberMux.Lock()
			if !unbounded {
				pendingRequests = pendingRequests - 1
			}
			subscriber.OnNext(t)
			subscriberMux.Unlock()
		},
		onComplete: func() {
			subscriberMux.Lock()
			mux.Lock()
			if publishersComplete {
				mux.Unlock()
				subscriber.OnComplete()
			} else {
				mux.Unlock()
				publishersSubscription.Request(1)
			}
			subscriberMux.Unlock()
		},
		onError: func(err error) {
			subscriberMux.Lock()
			publishersSubscription.Cancel()
			subscribtionMux.Lock()
			subscription.Cancel()
			subscribtionMux.Unlock()
			subscriber.OnError(err)
			subscriberMux.Unlock()
		},
	}

	return proc
}

func DistinctUntilChangedProcessor() cesium.Processor {
	var subscriber cesium.Subscriber
	var subscription cesium.Subscription
	subscriberMux := sync.Mutex{}
	subscribtionMux := sync.Mutex{}

	started := false
	var item cesium.T

	return &processor{
		subscribe: func(s cesium.Subscriber) cesium.Subscription {
			sub := &Subscription{
				CancelFunc: func() {
					subscribtionMux.Lock()
					if subscription != nil {
						subscription.Cancel()
					}
					subscribtionMux.Unlock()
				},
				RequestFunc: func(n int64) {
					subscribtionMux.Lock()
					if subscription != nil {
						subscription.Request(n)
					}
					subscribtionMux.Unlock()
				},
			}

			subscriberMux.Lock()
			subscriber = s
			subscriber.OnSubscribe(subscription)
			subscriberMux.Unlock()

			return sub
		},
		onSubscribe: func(s cesium.Subscription) {
			subscribtionMux.Lock()
			subscription = s
			subscribtionMux.Unlock()
		},
		onNext: func(t cesium.T) {
			if started  && item == t{
				subscribtionMux.Lock()
				subscription.Request(1)
				subscribtionMux.Unlock()
				return
			}

			started = true
			item = t
			subscriberMux.Lock()
			subscriber.OnNext(t)
			subscriberMux.Unlock()
		},
		onComplete: func() {
			subscriberMux.Lock()
			subscriber.OnComplete()
			subscriberMux.Unlock()
		},
		onError: func(err error) {
			subscriberMux.Lock()
			subscriber.OnError(err)
			subscriberMux.Unlock()
		},
	}
}

func TakeProcessor(n int64) cesium.Processor {
	var subscriber cesium.Subscriber
	var subscription cesium.Subscription
	subscriberMux := sync.Mutex{}
	subscribtionMux := sync.Mutex{}

	taken := int64(0)

	return &processor{
		subscribe: func(s cesium.Subscriber) cesium.Subscription {
			sub := &Subscription{
				CancelFunc: func() {
					subscribtionMux.Lock()
					if subscription != nil {
						subscription.Cancel()
					}
					subscribtionMux.Unlock()
				},
				RequestFunc: func(n int64) {
					subscribtionMux.Lock()
					if subscription != nil {
						subscription.Request(n)
					}
					subscribtionMux.Unlock()
				},
			}

			subscriberMux.Lock()
			subscriber = s
			subscriber.OnSubscribe(subscription)
			subscriberMux.Unlock()

			return sub
		},
		onSubscribe: func(s cesium.Subscription) {
			subscribtionMux.Lock()
			subscription = s
			subscribtionMux.Unlock()
		},
		onNext: func(t cesium.T) {
			if taken >= n {
				subscriberMux.Lock()
				subscriber.OnComplete()
				subscriberMux.Unlock()

				subscribtionMux.Lock()
				subscription.Cancel()
				subscribtionMux.Unlock()
				return
			}
			taken++
			subscriberMux.Lock()
			subscriber.OnNext(t)
			if taken >= n {
				subscriber.OnComplete()
				subscriberMux.Unlock()

				subscribtionMux.Lock()
				subscription.Cancel()
				subscribtionMux.Unlock()
				return
			}
			subscriberMux.Unlock()
		},
		onComplete: func() {
			subscriberMux.Lock()
			subscriber.OnComplete()
			subscriberMux.Unlock()
		},
		onError: func(err error) {
			subscriberMux.Lock()
			subscriber.OnError(err)
			subscriberMux.Unlock()
		},
	}
}