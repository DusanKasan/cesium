package internal

import (
	"sync"

	"github.com/DusanKasan/cesium"
)

// Just creates new cesium.Mono that emits the supplied item and completes.
func MonoJust(t cesium.T) cesium.Mono {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		if scheduler == nil {
			scheduler = SeparateGoroutineScheduler()
		}

		mux := sync.Mutex{}
		requested := false

		cancellable := scheduler.Schedule(func(canceller cesium.Canceller) {
			for {
				if canceller.IsCancelled() {
					return
				}

				mux.Lock()
				if requested {
					mux.Unlock()

					if canceller.IsCancelled() {
						return
					}

					subscriber.OnNext(t)
					break
				} else {
					mux.Unlock()
				}
			}

			if canceller.IsCancelled() {
				return
			}

			subscriber.OnComplete()
			return
		})

		return &Subscription{
			CancelFunc: func() {
				cancellable.Cancel()
			},
			RequestFunc: func(n int64) {
				mux.Lock()
				requested = true
				mux.Unlock()
			},
		}
	}

	return &Mono{OnSubscribe: onPublish}
}

// JustOrEmpty creates new cesium.Mono that emits the supplied item if it's non-nil
// and completes, otherwise just completes.
func MonoJustOrEmpty(t cesium.T) cesium.Mono {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		if scheduler == nil {
			scheduler = SeparateGoroutineScheduler()
		}

		mux := sync.Mutex{}
		requested := false

		cancellable := scheduler.Schedule(func(canceller cesium.Canceller) {
			if t != nil {
				for {
					if canceller.IsCancelled() {
						return
					}

					mux.Lock()
					if requested {
						mux.Unlock()
						if canceller.IsCancelled() {
							return
						}
						subscriber.OnNext(t)
						break
					} else {
						mux.Unlock()
					}
				}
			}

			if canceller.IsCancelled() {
				return
			}

			subscriber.OnComplete()
		})

		return &Subscription{
			CancelFunc: func() {
				cancellable.Cancel()
			},
			RequestFunc: func(n int64) {
				mux.Lock()
				requested = true
				mux.Unlock()
			},
		}
	}

	return &Mono{OnSubscribe: onPublish}
}

// FromCallable creates new cesium.Mono that emits the item returned from the supplied function. If the function
// returns nil, the returned Mono completes empty.
func MonoFromCallable(f func() cesium.T) cesium.Mono {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		if scheduler == nil {
			scheduler = SeparateGoroutineScheduler()
		}

		mux := sync.Mutex{}
		requested := false

		cancellable := scheduler.Schedule(func(canceller cesium.Canceller) {
			t := f()

			if t != nil {
				for {
					if canceller.IsCancelled() {
						return
					}

					mux.Lock()
					if requested {
						mux.Unlock()
						subscriber.OnNext(t)
						break
					} else {
						mux.Unlock()
					}
				}
			}

			if canceller.IsCancelled() {
				return
			}

			subscriber.OnComplete()

		})

		return &Subscription{
			CancelFunc: func() {
				cancellable.Cancel()
			},
			RequestFunc: func(n int64) {
				mux.Lock()
				requested = true
				mux.Unlock()
			},
		}
	}

	return &Mono{OnSubscribe: onPublish}
}

// Empty creates new cesium.Mono that emits no items and completes normally.
func MonoEmpty() cesium.Mono {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		if scheduler == nil {
			scheduler = SeparateGoroutineScheduler()
		}

		cancellable := scheduler.Schedule(func(canceller cesium.Canceller) {
			for {
				if canceller.IsCancelled() {
					return
				}

				subscriber.OnComplete()
				return
			}
		})

		return &Subscription{
			CancelFunc: func() {
				cancellable.Cancel()
			},
			RequestFunc: func(n int64) {
			},
		}
	}

	return &Mono{OnSubscribe: onPublish}
}

// Empty creates new cesium.Mono that emits no items and completes with error.
func MonoError(err error) cesium.Mono {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		if scheduler == nil {
			scheduler = SeparateGoroutineScheduler()
		}

		cancellable := scheduler.Schedule(func(canceller cesium.Canceller) {
			for {
				if canceller.IsCancelled() {
					return
				}

				subscriber.OnError(err)
				return
			}
		})

		return &Subscription{
			CancelFunc: func() {
				cancellable.Cancel()
			},
			RequestFunc: func(n int64) {
			},
		}
	}

	return &Mono{OnSubscribe: onPublish}
}

// Never creates new cesium.Mono that emits no items and never completes.
func MonoNever() cesium.Mono {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		if scheduler == nil {
			scheduler = SeparateGoroutineScheduler()
		}

		return &Subscription{
			CancelFunc: func() {
			},
			RequestFunc: func(n int64) {
			},
		}
	}

	return &Mono{OnSubscribe: onPublish}
}

// Defer creates new cesium.Mono by subscribing to the Mono returned from the supplied
// factory function for each subscribtion.
func MonoDefer(f func() cesium.Mono) cesium.Mono {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		if scheduler == nil {
			scheduler = SeparateGoroutineScheduler()
		}

		var subscription cesium.Subscription
		subscriptionMux := sync.Mutex{}

		canc := scheduler.Schedule(func(c cesium.Canceller) {
			subscriptionMux.Lock()
			subscription = f().Subscribe(subscriber)
			subscriptionMux.Unlock()
		})

		for {
			subscriptionMux.Lock()
			if subscription != nil {
				subscriptionMux.Unlock()
				break
			}
			subscriptionMux.Unlock()
		}

		return &Subscription{
			CancelFunc: func() {
				subscription.Cancel()
				canc.Cancel()
			},
			RequestFunc: func(n int64) {
				subscription.Request(n)
			},
		}
	}

	return &Mono{OnSubscribe: onPublish}
}

// Using Uses a resource, generated by a supplier for each individual Subscriber,
// while streaming the value from a Mono derived from the same resource and makes
// sure the resource is released if the sequence terminates or the Subscriber cancels.
func MonoUsing(resourceSupplier func() cesium.T, sourceSupplier func(cesium.T) cesium.Mono, resourceCleanup func(cesium.T)) cesium.Mono {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		if scheduler == nil {
			scheduler = SeparateGoroutineScheduler()
		}

		var subscription cesium.Subscription
		subscriptionMux := sync.Mutex{}
		resource := resourceSupplier()

		canc := scheduler.Schedule(func(c cesium.Canceller) {
			subscriptionMux.Lock()
			p := DoFinallyProcessor(func() {
				resourceCleanup(resource)
			})

			subscription = p.Subscribe(subscriber)
			subscriber.OnSubscribe(subscription)
			secondarySubscription := sourceSupplier(resource).Subscribe(p)
			p.OnSubscribe(secondarySubscription)

			subscriptionMux.Unlock()
		})

		for {
			subscriptionMux.Lock()
			if subscription != nil {
				subscriptionMux.Unlock()
				break
			}
			subscriptionMux.Unlock()
		}

		return &Subscription{
			CancelFunc: func() {
				subscription.Cancel()
				canc.Cancel()
			},
			RequestFunc: func(n int64) {
				subscription.Request(n)
			},
		}
	}

	return &Mono{OnSubscribe: onPublish}
}

// Create allows you to programmatically create a cesium.Mono with the
// capability of emitting multiple elements in a synchronous or asynchronous
// manner through the mono.Sink API
func MonoCreate(f func(cesium.MonoSink)) cesium.Mono {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		if scheduler == nil {
			scheduler = SeparateGoroutineScheduler()
		}

		sinkMux := sync.Mutex{}
		var sink *MonoSink

		cancellable := scheduler.Schedule(func(c cesium.Canceller) {
			sinkMux.Lock()
			sink = BufferMonoSink(subscriber, c)
			sinkMux.Unlock()
			f(sink)
		})

		for {
			sinkMux.Lock()
			if sink != nil {
				sinkMux.Unlock()
				break
			}
			sinkMux.Unlock()
		}

		return &Subscription{
			CancelFunc: func() {
				cancellable.Cancel()
			},
			RequestFunc: func(n int64) {
				sink.Request(n)
			},
		}
	}

	return &Mono{OnSubscribe: onPublish}
}
