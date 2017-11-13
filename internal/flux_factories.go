package internal

import (
	"math"
	"sync"

	"time"

	"github.com/DusanKasan/cesium"
)

type OverflowStrategy string

const OverflowStrategyBuffer = OverflowStrategy("buffer")
const OverflowStrategyDrop = OverflowStrategy("drop")
const OverflowStrategyError = OverflowStrategy("error")
const OverflowStrategyIgnore = OverflowStrategy("ignore")

// FluxJust creates new cesium.Flux that emits the supplied items.
func FluxJust(items ...cesium.T) cesium.Flux {
	return FluxFromSlice(items)
}

// fluxFromCallable is internal way to instantiate a ScalarFlux (a Flux that
// behaves like a mono, and offers more efficient implementation of some operators.
func fluxFromCallable(f func() (cesium.T, bool)) cesium.Flux {
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

					t, ok := f()
					if ok {
						subscriber.OnNext(t)
					}
					subscriber.OnComplete()
					return
				} else {
					mux.Unlock()
				}
			}
		})

		sub := &Subscription{
			CancelFunc: func() {
				cancellable.Cancel()
			},
			RequestFunc: func(n int64) {
				mux.Lock()
				requested = true
				mux.Unlock()
			},
		}

		subscriber.OnSubscribe(sub)
		return sub
	}

	return &ScalarFlux{
		&Flux{OnSubscribe: onPublish},
		f,
	}
}

// FromSlice creates new cesium.Flux that emits items from the supplied slice.
func FluxFromSlice(items []cesium.T) cesium.Flux {
	if len(items) == 0 {
		return FluxEmpty()
	}

	if len(items) == 1 {
		return fluxFromCallable(func() (cesium.T, bool) {
			return items[0], true
		})
	}

	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		if scheduler == nil {
			scheduler = SeparateGoroutineScheduler()
		}

		mux := sync.Mutex{}
		index := 0
		requested := int64(0)
		unbounded := false

		switch s := subscriber.(type) {
		case ConditionalSubscriber:
			fastPath := func(canceller cesium.Canceller) {
				for ; index < len(items); index++ {
					if canceller.IsCancelled() {
						return
					}

					s.OnNextIf(items[index])
				}

				if canceller.IsCancelled() {
					return
				}

				s.OnComplete()
			}

			cancellable := scheduler.Schedule(func(canceller cesium.Canceller) {
			sliceFor:
				for ; index < len(items); index++ {
					for {
						if canceller.IsCancelled() {
							return
						}

						mux.Lock()
						if unbounded {
							mux.Unlock()
							fastPath(canceller)
							return
						}
						mux.Unlock()

						mux.Lock()
						if requested > 0 {
							requested = requested - 1
							mux.Unlock()

							if canceller.IsCancelled() {
								return
							}

							if !s.OnNextIf(items[index]) {
								mux.Lock()
								requested = requested + 1
								mux.Unlock()
							}
							continue sliceFor
						} else {
							mux.Unlock()
						}
					}
				}

				if canceller.IsCancelled() {
					return
				}

				s.OnComplete()
			})

			sub := &Subscription{
				CancelFunc: func() {
					cancellable.Cancel()
				},
				RequestFunc: func(n int64) {
					mux.Lock()
					if unbounded {
						mux.Unlock()
						return
					}

					if n == math.MaxInt64 {
						unbounded = true
					} else {
						requested = requested + n
					}
					mux.Unlock()
				},
			}

			s.OnSubscribe(sub)
			return sub
		default:
			fastPath := func(canceller cesium.Canceller) {
				for ; index < len(items); index++ {
					if canceller.IsCancelled() {
						return
					}

					subscriber.OnNext(items[index])
				}

				if canceller.IsCancelled() {
					return
				}

				subscriber.OnComplete()
			}

			cancellable := scheduler.Schedule(func(canceller cesium.Canceller) {
			sliceFor:
				for ; index < len(items); index++ {
					for {
						if canceller.IsCancelled() {
							return
						}

						mux.Lock()
						if unbounded {
							mux.Unlock()
							fastPath(canceller)
							return
						}
						mux.Unlock()

						mux.Lock()
						if requested > 0 {
							requested = requested - 1
							mux.Unlock()

							if canceller.IsCancelled() {
								return
							}

							subscriber.OnNext(items[index])
							continue sliceFor
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

			sub := &Subscription{
				CancelFunc: func() {
					cancellable.Cancel()
				},
				RequestFunc: func(n int64) {
					mux.Lock()
					if unbounded {
						mux.Unlock()
						return
					}

					if n == math.MaxInt64 {
						unbounded = true
					} else {
						requested = requested + n
					}
					mux.Unlock()
				},
			}

			subscriber.OnSubscribe(sub)
			return sub
		}
	}

	return &Flux{OnSubscribe: onPublish}
}

// Range creates new cesium.Flux that emits 64bit integers from start to
// (start + count).
func FluxRange(start int, count int) cesium.Flux {
	if count == 0 {
		return FluxEmpty()
	}

	if count == 1 {
		return FluxJust(start)
	}

	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		if scheduler == nil {
			scheduler = SeparateGoroutineScheduler()
		}

		mux := sync.Mutex{}
		requested := int64(0)
		unbounded := false
		item := int64(start)

		switch s := subscriber.(type) {
		case ConditionalSubscriber:
			cancellable := scheduler.Schedule(func(canceller cesium.Canceller) {
			sliceFor:
				for ; item < int64(start)+int64(count); item++ {
					for {
						if canceller.IsCancelled() {
							return
						}

						mux.Lock()
						if unbounded {
							mux.Unlock()
							for ; item < int64(start)+int64(count); item++ {
								if canceller.IsCancelled() {
									return
								}

								s.OnNextIf(item)
							}

							s.OnComplete()
							return
						}
						mux.Unlock()

						mux.Lock()
						if requested > 0 {
							requested = requested - 1
							mux.Unlock()
							if canceller.IsCancelled() {
								return
							}

							if !s.OnNextIf(item) {
								mux.Lock()
								requested = requested + 1
								mux.Unlock()
							}

							continue sliceFor
						} else {
							mux.Unlock()
						}
					}
				}

				if canceller.IsCancelled() {
					return
				}

				s.OnComplete()
			})

			sub := &Subscription{
				CancelFunc: func() {
					cancellable.Cancel()
				},
				RequestFunc: func(n int64) {
					mux.Lock()
					if unbounded {
						mux.Unlock()
						return
					}

					if n == math.MaxInt64 {
						unbounded = true
					} else {
						requested = requested + n
					}
					mux.Unlock()
				},
			}

			s.OnSubscribe(sub)
			return sub
		default:
			cancellable := scheduler.Schedule(func(canceller cesium.Canceller) {
			sliceFor:
				for ; item < int64(start)+int64(count); item++ {
					for {
						if canceller.IsCancelled() {
							return
						}

						mux.Lock()
						if unbounded {
							mux.Unlock()
							for ; item < int64(start)+int64(count); item++ {
								if canceller.IsCancelled() {
									return
								}

								subscriber.OnNext(item)
							}

							subscriber.OnComplete()
							return
						}
						mux.Unlock()

						mux.Lock()
						if requested > 0 {
							requested = requested - 1
							mux.Unlock()
							if canceller.IsCancelled() {
								return
							}
							subscriber.OnNext(item)
							continue sliceFor
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

			sub := &Subscription{
				CancelFunc: func() {
					cancellable.Cancel()
				},
				RequestFunc: func(n int64) {
					mux.Lock()
					if unbounded {
						mux.Unlock()
						return
					}

					if n == math.MaxInt64 {
						unbounded = true
					} else {
						requested = requested + n
					}
					mux.Unlock()
				},
			}

			subscriber.OnSubscribe(sub)
			return sub
		}

	}

	return &Flux{OnSubscribe: onPublish}
}

// Empty creates new cesium.Flux that emits no items and completes normally.
func FluxEmpty() cesium.Flux {
	return fluxFromCallable(func() (cesium.T, bool) {
		return nil, false
	})
}

// Error creates new cesium.Flux that emits no items and completes with error.
func FluxError(err error) cesium.Flux {
	//TODO: error flux

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

		sub := &Subscription{
			CancelFunc: func() {
				cancellable.Cancel()
			},
			RequestFunc: func(n int64) {
			},
		}

		subscriber.OnSubscribe(sub)
		return sub
	}

	return &Flux{OnSubscribe: onPublish}

}

// Never creates new cesium.Flux that emits no items and never completes.
func FluxNever() cesium.Flux {
	//TODO: Noop flux

	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		if scheduler == nil {
			scheduler = SeparateGoroutineScheduler()
		}

		sub := &Subscription{
			CancelFunc: func() {
			},
			RequestFunc: func(n int64) {
			},
		}

		subscriber.OnSubscribe(sub)
		return sub
	}

	return &Flux{OnSubscribe: onPublish}

}

// Defer creates new cesium.Flux by subscribing to the Publisher returned from
// the supplied factory function for each subscribtion.
func FluxDefer(f func() cesium.Publisher) cesium.Flux {
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

		sub := &Subscription{
			CancelFunc: func() {
				subscription.Cancel()
				canc.Cancel()
			},
			RequestFunc: func(n int64) {
				subscription.Request(n)
			},
		}

		subscriber.OnSubscribe(sub)
		return sub
	}

	return &Flux{OnSubscribe: onPublish}
}

// Using uses a resource, generated by a supplier for each individual Subscriber,
// while streaming the value from a Publisher derived from the same resource
// and makes sure the resource is released if the sequence terminates or the
// Subscriber cancels.
func FluxUsing(resourceSupplier func() cesium.T, sourceSupplier func(cesium.T) cesium.Publisher, resourceCleanup func(cesium.T)) cesium.Flux {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		if scheduler == nil {
			scheduler = SeparateGoroutineScheduler()
		}

		subscription := &BufferedProxySubscription{}
		resource := resourceSupplier()

		canc := scheduler.Schedule(func(c cesium.Canceller) {
			p := DoFinallyProcessor(func() {
				resourceCleanup(resource)
			})

			s := p.Subscribe(DoObserver(
				func(t cesium.T) {
					subscriber.OnNext(t)
				},
				func() {
					subscriber.OnComplete()
				},
				func(e error) {
					subscriber.OnError(e)
				},
			))

			secondarySubscription := sourceSupplier(resource).Subscribe(p)
			p.OnSubscribe(secondarySubscription)

			subscription.SetSubscription(s)
		})

		sub := &Subscription{
			CancelFunc: func() {
				subscription.Cancel()
				canc.Cancel()
			},
			RequestFunc: func(n int64) {
				subscription.Request(n)
			},
		}

		subscriber.OnSubscribe(sub)
		return sub
	}

	return &Flux{OnSubscribe: onPublish}

}

// Create allows you to programmatically create a cesium.Flux with the
// capability of emitting multiple elements in a synchronous or asynchronous
// manner through the flux.Sink API
func FluxCreate(f func(cesium.FluxSink), overflowStrategy OverflowStrategy) cesium.Flux {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		if scheduler == nil {
			scheduler = SeparateGoroutineScheduler()
		}

		var sink *FluxSink
		sinkMux := sync.Mutex{}

		cancellable := scheduler.Schedule(func(c cesium.Canceller) {
			sinkMux.Lock()
			switch overflowStrategy {
			case OverflowStrategyBuffer:
				sink = BufferFluxSink(subscriber, c)
			case OverflowStrategyDrop:
				sink = DropFluxSink(subscriber, c)
			case OverflowStrategyError:
				sink = ErrorFluxSink(subscriber, c)
			case OverflowStrategyIgnore:
				sink = IgnoreFluxSink(subscriber, c)
			}
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

		sub := &Subscription{
			CancelFunc: func() {
				cancellable.Cancel()
			},
			RequestFunc: func(n int64) {
				sink.Request(n)
			},
		}

		subscriber.OnSubscribe(sub)
		return sub
	}

	return &Flux{OnSubscribe: onPublish}
}

func FluxGenerate(f func(cesium.SynchronousSink)) cesium.Flux {
	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		if scheduler == nil {
			scheduler = SeparateGoroutineScheduler()
		}

		requested := int64(0)
		requestedMux := sync.Mutex{}
		unbounded := false

		switch s := subscriber.(type) {
		case ConditionalSubscriber:
			cancellable := scheduler.Schedule(func(c cesium.Canceller) {
				for !c.IsCancelled() {
					requestedMux.Lock()
					if requested > 0 || unbounded {
						requested = requested - 1
						sink := &SynchronousSink{}

						f(sink)

						if !c.IsCancelled() {
							sig := sink.Signal()
							if sig == nil {
								s.OnError(cesium.NoEmissionOnSynchronousSinkError)
								requestedMux.Unlock()
								return
							}

							switch sig.Type() {
							case cesium.SignalTypeOnNext:
								if !s.OnNextIf(sig.Item()) {
									requested = requested + 1
								}
								requestedMux.Unlock()
							case cesium.SignalTypeOnComplete, cesium.SignalTypeOnError:
								sig.Accept(s)
								requestedMux.Unlock()
								return
							}
						}
					}
					requestedMux.Unlock()
				}
			})

			sub := &Subscription{
				CancelFunc: func() {
					cancellable.Cancel()
				},
				RequestFunc: func(n int64) {
					requestedMux.Lock()
					if unbounded {
						requestedMux.Unlock()
						return
					}

					if n == math.MaxInt64 {
						unbounded = true
					}

					requested = requested + n
					requestedMux.Unlock()
				},
			}

			s.OnSubscribe(sub)
			return sub
		default:
			cancellable := scheduler.Schedule(func(c cesium.Canceller) {
				for !c.IsCancelled() {
					requestedMux.Lock()
					if requested > 0 || unbounded {
						requested = requested - 1
						sink := &SynchronousSink{}

						f(sink)

						if !c.IsCancelled() {
							sig := sink.Signal()
							if sig == nil {
								subscriber.OnError(cesium.NoEmissionOnSynchronousSinkError)
								requestedMux.Unlock()
								return
							}

							sig.Accept(subscriber)

							if sig.IsTerminal() {
								requestedMux.Unlock()
								return
							}
						}
					}
					requestedMux.Unlock()
				}
			})

			sub := &Subscription{
				CancelFunc: func() {
					cancellable.Cancel()
				},
				RequestFunc: func(n int64) {
					requestedMux.Lock()
					if unbounded {
						requestedMux.Unlock()
						return
					}

					if n == math.MaxInt64 {
						unbounded = true
					}

					requested = requested + n
					requestedMux.Unlock()
				},
			}

			subscriber.OnSubscribe(sub)
			return sub
		}
	}

	return &Flux{OnSubscribe: onPublish}
}

func FluxFromChannel(ch <-chan cesium.T) cesium.Flux {
	subscribed := false
	mux := sync.Mutex{}

	onPublish := func(subscriber cesium.Subscriber, scheduler cesium.Scheduler) cesium.Subscription {
		mux.Lock()
		if subscribed {
			panic("can not resubscribe to flux made from a channel")
		}

		subscribed = true
		mux.Unlock()

		if scheduler == nil {
			scheduler = SeparateGoroutineScheduler()
		}

		requestedMux := sync.Mutex{}
		requested := int64(0)
		unbounded := false

		canc := scheduler.Schedule(func(c cesium.Canceller) {
			if c.IsCancelled() {
				return
			}

			for {
				requestedMux.Lock()
				if requested == 0 && !unbounded {
					requestedMux.Unlock()
					continue
				}
				requestedMux.Unlock()

				select {
				case t, ok := <-ch:
					if c.IsCancelled() {
						return
					}

					if ok {
						requestedMux.Lock()
						requested--
						requestedMux.Unlock()
						subscriber.OnNext(t)
					} else {
						subscriber.OnComplete()
						return
					}
				case <-time.After(time.Millisecond):
					if c.IsCancelled() {
						return
					}
				}
			}
		})

		sub := &Subscription{
			CancelFunc: func() {
				canc.Cancel()
			},
			RequestFunc: func(n int64) {
				requestedMux.Lock()
				if unbounded {
					requestedMux.Unlock()
					return
				}

				requested = requested + n
				requestedMux.Unlock()
			},
		}

		subscriber.OnSubscribe(sub)
		return sub
	}

	return &Flux{OnSubscribe: onPublish}
}
