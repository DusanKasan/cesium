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

type conditionalProcessor struct {
	*processor
	onNextIf func(cesium.T) bool
}

func (c *conditionalProcessor) OnNextIf(t cesium.T) bool {
	return c.onNextIf(t)
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

	p := &processor{
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
						subscribtionMux.Unlock()
						subscription.Request(n)
						return
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

	return &conditionalProcessor{
		p,
		func(t cesium.T) bool {
			b := f(t)

			if b {
				switch sub := subscriber.(type) {
				case ConditionalSubscriber:
					// Propagate the microfusion downstream
					subscriberMux.Lock()
					x := sub.OnNextIf(t)
					subscriberMux.Unlock()
					return x
				default:
					subscriberMux.Lock()
					subscriber.OnNext(t)
					subscriberMux.Unlock()
				}
			}

			return b
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
			subscriber.OnSubscribe(sub)
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
			if terminated {
				subscriberMux.Unlock()
				return
			}
			subscriber.OnComplete()
			subscriberMux.Unlock()
		},
		onError: func(err error) {
			subscriberMux.Lock()
			if terminated {
				subscriberMux.Unlock()
				return
			}
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
			if started && item == t {
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

type indexedEmission struct {
	t     cesium.T
	index int
}

func FlatMapProcessor(f func(cesium.T) cesium.Publisher, scheduler cesium.Scheduler) cesium.Processor {
	var subscriber cesium.Subscriber
	var subscription cesium.Subscription
	subscriberMux := sync.Mutex{}
	subscribtionMux := sync.Mutex{}

	var emissionBuffer []indexedEmission
	subscriptions := make(map[int]cesium.Subscription)
	mux := sync.Mutex{}
	hasItems := false
	closed := false
	currentIndex := 0
	openSubscriptions := 0
	requested := int64(0)
	mainEmittedAll := false

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
					mux.Lock()
					for len(emissionBuffer) > 0 && n > 0 {
						emission := emissionBuffer[0]
						emissionBuffer = emissionBuffer[1:]
						n--

						scheduler.Schedule(func(c cesium.Canceller) {
							if !c.IsCancelled() {
								subscriber.OnNext(emission.t)
								mux.Lock()
								subscriptions[emission.index].Request(1)
								mux.Unlock()
							}
						})
					}

					if openSubscriptions == 0 && mainEmittedAll && len(emissionBuffer) == 0 {
						scheduler.Schedule(func(c cesium.Canceller) {
							subscriber.OnComplete()
						})
					}

					requested = requested + n
					mux.Unlock()
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
			s.Request(math.MaxInt64)
			subscribtionMux.Unlock()
		},
		onNext: func(t cesium.T) {
			mux.Lock()
			if closed {
				mux.Unlock()
				return
			}

			hasItems = true
			i := currentIndex
			currentIndex++
			openSubscriptions++
			mux.Unlock()

			sub := f(t).Subscribe(DoObserver(
				func(t cesium.T) {
					mux.Lock()
					if !closed {
						emissionBuffer = append(emissionBuffer, indexedEmission{t, i})
						if requested > 0 {
							requested--
							emission := emissionBuffer[0]
							emissionBuffer = emissionBuffer[1:]
							subscriptions[i].Request(1)
							mux.Unlock()

							scheduler.Schedule(func(c cesium.Canceller) {
								if !c.IsCancelled() {
									subscriberMux.Lock()
									subscriber.OnNext(emission.t)
									subscriberMux.Unlock()
								}
							})

							return
						}

						mux.Unlock()
						return
					}
					mux.Unlock()
				},
				func() {
					mux.Lock()
					openSubscriptions--
					if openSubscriptions == 0 && mainEmittedAll && len(emissionBuffer) == 0 {
						scheduler.Schedule(func(c cesium.Canceller) {
							if !c.IsCancelled() {
								subscriber.OnComplete()
							}
						})

					}
					mux.Unlock()
				},
				func(err error) {
					mux.Lock()
					closed = true
					emissionBuffer = []indexedEmission{}

					for _, s := range subscriptions {
						s.Cancel()
					}
					mux.Unlock()

					subscriberMux.Lock()
					subscriber.OnError(err)
					subscriberMux.Unlock()
				},
			))

			mux.Lock()
			subscriptions[i] = sub
			mux.Unlock()
			sub.Request(1)
		},
		onComplete: func() {
			mux.Lock()
			mainEmittedAll = true
			mux.Unlock()

			subscriberMux.Lock()
			if !hasItems {
				subscriber.OnComplete()
			}
			subscriberMux.Unlock()
		},
		onError: func(err error) {
			mux.Lock()
			closed = true
			emissionBuffer = []indexedEmission{}

			for _, s := range subscriptions {
				s.Cancel()
			}
			mux.Unlock()

			subscriberMux.Lock()
			subscriber.OnError(err)
			subscriberMux.Unlock()
		},
	}
}
