// Cesium is a general purpose 4th generation  non-blocking reactive library
// that offers efficient demand management (in the form of managing
// "backpressure").  It offers composable asynchronous sequence APIs Flux (for
// [N] elements) and Mono (for [0|1] elements), extensively implementing the
// Reactive Extensions specification where possible.
//
// Cesium also ships with a testing framework under cesium/verifier that allows
// you to test your reactive code easily.
//
// Under cesium/commons you'll find various bindings for reactive encoding/decoding
// and reactive network engine (wip).
package cesium

import (
	"log"
	"time"
)

// T is a type placeholder that cesium uses almost everywhere to mask the lack
// of generics in Go.
type T interface{}

// Publisher is Observable that supports asynchronous backpressure. This is
// implemented via the Request(int) method of a Subscription returned from the
// Subscribe method which will instruct the Publisher to emit requested amount
// of items to the Subscriber. Publishers come in two varieties, hot and cold
// publishers. Cold publishers will replay all items for each subscription
// (think slices) while hot ones will only emit new items (like channels).
//
// The hot publishers returned by cesium (think FromChannel) only support one
// subscriber and subscribing multiple times returns empty subscription.
//
// It is important to note, that if you are implementing publisher you must
// handle the fact that different Subscriber's Subscriptions should not
// interfere with each other. For example requesting from one subscription
// shouldn't emit items to others. This is hard to achieve for hot Publishers
// as it has to be implemented with a per subscription buffer inside the
// publisher. This is why cesium does the trade-off of supporting only one
// subscriber for hot observables.
type Publisher interface {
	// Subscribe will subscribe the passed Subscriber to this Publisher and
	// returns a Subscription that can also be used to control the Publisher or
	// cancel the subscription. In the process of subscribing the
	// OnSubscription(Subscription) method ofSubscriber will be invoked.
	Subscribe(Subscriber) Subscription
}

type Subscriber interface {
	// OnNext is called by a Publisher when a it emits an item T.
	OnNext(T)

	// OnError is called by a Publisher when a it completes with an error.
	OnError(error)

	// OnComplete is called by a Publisher when a it completes successfully.
	OnComplete()

	// OnSubscribe is called by a Publisher when a this subscriber subscribes to
	// it.
	OnSubscribe(Subscription)
}

type Subscription interface {
	// Cancel the Subscription, effectively telling the Publisher to stop
	// emitting on this subscription.
	Cancel()

	// Request requqests the specified number of items from the Publisher.
	Request(int64)

	// RequestUnbounded switches the Publisher to unbounded mode where it acts
	// like a observable, emitting items without request. Internally this is
	// equivalent to calling Request(math.MaxInt64).
	RequestUnbounded()
}

// Processor is a Publisher that is also a Subscriber, effectively making it a
// reactive operator.
type Processor interface {
	Subscriber
	Publisher
}

// Observables/Publishers emit on the thread Subscribe was called on, so
// to introduce multi-threading we execute everything on schedulers. Some
// operators allow you to pass a specific scheduler, because they can not
// re-emit items on the same scheduler they received them.
type Scheduler interface {
	// Schedulers must supply Canceller to the action and call Cancel in
	// Cancel method of the returned Cancellable. This is done like this
	// because there is no way to kill a goroutine from the outside.
	Schedule(action func(Canceller)) Cancellable
}

// Cancellable is a way to cancel an action scheduled on a Scheduler.
type Cancellable interface {
	// Cancel the scheduled action.
	Cancel()
}

// Canceller is used to propagate cancellation to long running scheduled tasks,
// as goroutines cannot be cancelled from the outside.
type Canceller interface {
	// IsCancelled checks if the running action was already requested to be
	// cancelled. If this returns true, the action should be cancelled.
	IsCancelled() bool

	// Register an callback that will be executed when the action is cancelled.
	OnCancel(func())
}

// Flux is a publisher with reactive operators that emits 0 to N elements, and
// then completes (successfully or with an error).
type Flux interface {
	Publisher

	Filter(func(T) bool) Flux
	DistinctUntilChanged() Flux
	Take(int64) Flux

	Map(func(T) T) Flux
	FlatMap(func(T) Publisher, ...Scheduler) Flux

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

	Handle(func(T, SynchronousSink)) Flux

	Count() Mono
	Reduce(func(T, T) T) Mono
	Scan(func(T, T) T) Flux
	All(func(T) bool) Mono
	Any(func(T) bool) Mono
	HasElements() Mono
	HasElement(T) Mono
	Concat(Publisher /*<cesium.Publisher>*/) Flux
	ConcatWith(...Publisher) Flux

	OnErrorReturn(T) Flux

	BlockFirst() (T, bool, error)
	BlockFirstTimeout(time.Duration) (T, bool, error)
	BlockLast() (T, bool, error)
	BlockLastTimeout(time.Duration) (T, bool, error)
}

// Mono is a publisher with reactive operators that emits 0 or 1 elements, and
// then completes (successfully or with an error).
type Mono interface {
	Publisher

	Map(func(T) T) Mono
	FlatMap(fn func(T) Mono, scheduler ...Scheduler) Mono
	FlatMapMany(fn func(T) Publisher, scheduler ...Scheduler) Flux

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

	Handle(func(T, SynchronousSink)) Mono

	ConcatWith(...Publisher) Flux

	OnErrorReturn(T) Mono

	Block() (T, bool, error)
	BlockTimeout(time.Duration) (T, bool, error)
}

// FluxSink is used in the flux.Create constructor to create a Flux
// programmatically.
type FluxSink interface {
	// Emit T.
	Next(T)

	// Emit complete signal.
	Complete()

	// Emit error signal.
	Error(error)

	// Check if cancellation war requested.
	IsCancelled() bool

	// Register a callback to be executed upon cancellation.
	OnCancel(func())

	// Register a callback to be executed when the returned Flux terminates by
	// completing (successfully or error) or by cancellation.
	OnDispose(func())

	// Register a callback to be executed when items are requesting from the
	// returned Flux.
	OnRequest(func(int64))

	// Returns the current outstanding request amount.
	RequestedFromDownstream() int64
}

// SynchronousSink is used in the flux/mono Generate factory functions as a
// means to generate the returned publisher emissions sequentially.
type SynchronousSink interface {
	// Emit T.
	Next(T)

	// Emit complete signal.
	Complete()

	// Emit error signal.
	Error(error)
}

// MonoSink is used in the mono.Create constructor to create a Mono
// programmatically.
type MonoSink interface {
	// Emit T and complete.
	CompleteWith(T)

	// Complete empty.
	Complete()

	// Emit the error signal.
	Error(error)

	// Register a callback to be executed upon cancellation.
	OnCancel(func())

	// Register a callback to be executed when the returned Flux terminates by
	// completing (successfully or error) or by cancellation.
	OnDispose(func())
}

// MaterializedEmission is used in (De)Materialize operator.
type MaterializedEmission struct {
	// The type of emission: "next", "complete", "error"
	EventType string

	// If the emission was of type "error" the error will be in this field
	Err error

	// If the emission was of type "next" the item will be in this field.
	Value T
}
