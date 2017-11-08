package cesium

import "log"

type T interface{}

// Publisher/Subscriber contract
//
// Publisher is Observable that supports asynchronous backpressure. This is
// implemented via the Request(int) method of a Subscription returned from the
// Subscribe method which will instruct the Publisher to emit requested amount
// of items to the Subscriber.

type Publisher interface {
	Subscribe(Subscriber) Subscription
}

type Subscriber interface {
	OnNext(T)
	OnError(error)
	OnComplete()
	OnSubscribe(Subscription)
}

type Cancellable interface {
	Cancel()
}

type Subscription interface {
	Cancellable
	Request(int64)
	RequestUnbounded()
}

type Processor interface {
	Subscriber
	Publisher
}

// Scheduler
//
// Observables/Publishers emit on the thread Subscribe was called on, so
// to introduce multi-threading we execute everything on schedulers. Some
// operators allow you to pass a specific scheduler, because they can not
// re-emit items on the same scheduler they received them.

type Scheduler interface {
	// Schedulers must supply Canceller to the action and call Cancel in
	// Cancel method of the returned cancellable. This is done like this
	// because there is no way to kill a goroutine from the outside.
	Schedule(action func(Canceller)) Cancellable
}

type Canceller interface {
	IsCancelled() bool
	OnCancel(func())
}

type Flux interface {
	Publisher

	Filter(func(T) bool) Flux
	DistinctUntilChanged() Flux
	Take(int64) Flux

	Map(func(T) T) Flux

	DoOnSubscribe(func(Subscription)) Flux
	DoOnRequest(func(int64)) Flux
	DoOnNext(func(T)) Flux
	DoOnError(func(error)) Flux
	DoOnComplete(func()) Flux
	DoOnTerminate(func()) Flux
	DoAfterTerminate(func()) Flux
	DoFinally(func()) Flux
	DoOnCancel(func()) Flux
	Log(*log.Logger) Flux

	Handle(fn func(T, SynchronousSink)) Flux

	Count() Mono
	Reduce(func(T, T) T) Mono
	Scan(func(T, T) T) Flux
	All(func(T) bool) Mono
	Any(func(T) bool) Mono
	HasElements() Mono
	HasElement(T) Mono
	Concat(Publisher /*<cesium.Publisher>*/) Flux
	ConcatWith(...Publisher) Flux
}

type Mono interface {
	Publisher

	Map(func(T) T) Mono

	Filter(func(T) bool) Mono

	DoOnSubscribe(func(Subscription)) Mono
	DoOnRequest(func(int64)) Mono
	DoOnNext(func(T)) Mono
	DoOnError(func(error)) Mono
	DoOnSuccess(func()) Mono
	DoOnTerminate(func()) Mono
	DoAfterTerminate(func()) Mono
	DoOnCancel(func()) Mono
	DoFinally(func()) Mono
	Log(*log.Logger) Mono

	ConcatWith(...Publisher) Flux
}

type FluxSink interface {
	Next(T)
	Complete()
	Error(error)

	IsCancelled() bool
	OnCancel(func())
	OnDispose(func())
	OnRequest(func(int64))

	RequestedFromDownstream() int64
}

type SynchronousSink interface {
	Next(T)
	Complete()
	Error(error)
}

type MonoSink interface {
	CompleteWith(T)
	Complete()
	Error(error)

	OnCancel(func())
	OnDispose(func())
}

type MaterializedEmission struct {
	EventType string // "next", "complete", "error"
	Err       error
	Value     T
}
