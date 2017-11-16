# Cesium

** This library is a work in progress **

[![Coverage Status](https://coveralls.io/repos/github/DusanKasan/cesium/badge.svg?branch=master)](https://coveralls.io/github/DusanKasan/cesium?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/DusanKasan/cesium)](https://goreportcard.com/report/github.com/DusanKasan/cesium) [![CircleCI](https://circleci.com/gh/DusanKasan/cesium.svg?style=shield)](https://circleci.com/gh/DusanKasan/cesium) [![GoDoc](https://godoc.org/github.com/DusanKasan/cesium?status.svg)](https://godoc.org/github.com/DusanKasan/cesium)

This is a port of [Project Reactor](https://projectcesium.io/) into [Go](https://golang.org/). It provides reactive data streams with asynchronous pull backpressure and operator fusion. Its aim is to be as close to the proposed Java API as possible, altering it slightly where needed for it to make sense in Go.

For more information see generated godocs about:
- [Operators and usage](https://godoc.org/github.com/DusanKasan/cesium)
- [Instantiating a Flux](https://godoc.org/github.com/DusanKasan/cesium/flux)
- [Instantiating a Mono](https://godoc.org/github.com/DusanKasan/cesium/mono)

More thorough documentation to come after most/all operators are implemented.

### Naming

The explanations contained here assume the knowledge of the [Observer Pattern](https://sourcemaking.com/design_patterns/observer).

- **Publisher** is an observable that supports asynchronous backpressure in form of its `Request(int64)` method. It will only ever emit the requested amount of emissions (does not apply to the closing emissions complete and error). It emission footprint is ( \[Next\](0-N) Complete|Error )
- **Flux** is a Publisher specific to this library that has access to operators described here. It can be switched to unbounded mode, basically transforming it to an observable by calling `Request(math.MaxInt64)` on its Subscripiton.
- **Mono** is a Publisher specific to this library whose emission footprint is ( \[Next\](0-1) Complete|Error ). It also has operators described here. It can be switched to unbounded mode, basically transforming it to an observable by calling `RequestUnbounded()` which is really just a proxy to `Request(math.MaxInt64)` on its Subscripiton.
- **Subscription** is the result of subscribing a Subscriber to a Publisher. Subscription provides means to control the pull-backpressure, via the `Request(int64)` method, that will intstruct the Publisher to emit specified amount of items. It also serves as a mean of emission cancellation via its `Cancel()` method that will cause the Publisher to stop emitting and shut down.
- **Subscriber** is an observer that will receive a Subscription object upon subscription. This is achieved via the Subscriber `OnSubscribe(Subscription)` method
- **Observer** in the scope of this library means a Publisher that does not control its subscription.
- **Scheduler** and its only method `Schedule(func(Canceller)) Cancellable`is where everything is executed. It allows us to execute different things on different threads. The returned Cancelable can be `Cancel()`ed. When cancelled the Canceller's method `IsCanceled()` will return true. This is a small hindrance and is done like this because Go goroutines can be only cancelled from inside.

### Simple examples

#### Subscribing an observer

Observer itself does not support backpressure so we have two options, apply an operator that will request emissions from the Publisher automatically or we can control this behavior via the returned Subscription object. In this example we do the latter.

```Go
subscription := flux.FromSlice(
    []cesium.T{1, 2, 3},
).Map(func(t cesium.T) cesium.T {
    return t.(int) + 1
}).Filter(func(t cesium.T) bool {
    return t.(int) < 4
}).Subscribe(PrintObserver())

subscription.Request(1)
// Printsubscriber will asynchronously print 2

subscription.Request(1)
// Printsubscriber will asynchronously print 3

subscription.Request(1)
// Printsubscriber will receive Complete signal
```

#### Subscribing an Subscriber

```Go
// This subscriber provides a Request(int64) method so we can control it
s := ControlledSubscriber()

flux.FromSlice(
    []cesium.T{1, 2, 3},
).Map(func(t cesium.T) cesium.T {
    return t.(int) + 1
}).Filter(func(t cesium.T) bool {
    return t.(int) < 3
}).Subscribe(s) // s will receive Subscription

s.Request(2) // s will receive 2 emissions
```

#### Switching to unbounded mode

```Go
subscription := flux.FromSlice(
    []cesium.T{1, 2, 3},
).Map(func(t cesium.T) cesium.T {
    return t.(int) + 1
}).Filter(func(t cesium.T) bool {
    return t.(int) < 4
}).Subscribe(PrintObserver())

subscription.RequestUnbounded() // Publisher will not wait for Request() to emit
```

### Operator implementation progress

Operators listed according to [Reactor docs](https://projectcesium.io/docs/core/release/reference/docs/index.html)

#### Factories
- [x] Just
- [x] Mono.JustOrEmpty
- [ ] Mono.FromSupplier
- [x] FromSlice
- [x] FromChannel
- [x] Mono.FromCallable
- [x] Empty
- [x] Never
- [x] Error
- [x] Defer
- [x] Using
- [x] Flux.Generate
- [x] Create
- [ ] Interval

#### Transforming
- [x] Map(func(T) T)
- [ ] Cast
- [x] FlatMap
- [x] Handle(func(T, SynchronousSink))
- [ ] Flux.FlatMapSequential
- [x] Mono.FlatMapMany
- [x] Flux.ToSlice
    - Maybe ToList (LinkedList would be better to handle large datasets)
- [ ] Flux.ToSortedSlice
- [ ] Flux.ToMap
- [x] Flux.ToChannel
- [x] Flux.Count()
- [x] Flux.Reduce(func(T, T) T)
- [x] Flux.Scan(func(T, T) T)
- [x] Flux.All(func(T) bool)
- [x] Flux.Any(func(T) bool)
- [x] Flux.HasElements()
- [x] Flux.HasElement(T) Flux
- [x] Flux.Concat(Publisher<Publisher>) Flux
- [x] ConcatWith(Publisher) Flux
- [ ] Flux.ConcatDelayError
- [ ] Flux.MergeSequential
- [ ] Flux.Merge
- [ ] MergeWith
- [ ] Zip
- [ ] ZipWith
- [ ] Mono.And
- [ ] Mono.When
- [ ] Flux.CombineLatest
- [ ] First (implement before Or)
- [ ] Or
- [ ] SwitchMap
- [ ] SwitchOnNext
- [ ] Repeat
- [ ] SwitchIfEmpty
- [ ] IgnoreElements
- [ ] Then
- [ ] ThenEmpty
- [ ] ThenMany
- [ ] Mono.DelayUntilOther
- [ ] Mono.DelayUntil
- [ ] Expand
- [ ] ExpandDeep

#### Peeking

- [x] DoOnNext(func(T))
- [x] Flux.DoOnComplete
- [x] Mono.DoOnSuccess
- [x] DoOnError(func(error))
- [x] DoOnCancel(func())
- [x] DoOnSubscribe(func(Subscription))
- [x] DoOnRequest
- [x] DoOnTerminate
- [x] DoAfterTerminate
- [x] DoFinally(func())
- [x] Log(log.Logger)
- [x] DoOnEach
- [x] Materialize
- [x] Dematerialize

#### Filtering

- [x] Filter
- [ ] FilterWhen
- [ ] OfType
- [ ] Flux.Distinct
- [x] Flux.DistinctUntilChanged
- [x] Flux.Take
- [ ] Flux.TakeInPeriod
- [ ] Flux.Next
- [ ] Flux.LimitRequest
- [ ] Flux.TakeUntil
- [ ] Flux.TakeUntilOther
- [ ] Flux.TakeWhile
- [ ] Flux.ElementAt
- [ ] Flux.TakeLast
- [ ] Flux.Last
- [ ] Flux.LastOrDefault
- [ ] Flux.LastOrDefault
- [ ] Flux.Skip
- [ ] Flux.SkipPeriod
- [ ] Flux.SkipLast
- [ ] Flux.SkipUntil
- [ ] Flux.SkipUntilOther
- [ ] Flux.SkipWhile
- [ ] Flux.Sample
- [ ] Flux.SampleFirst
- [ ] Flux.SampleUsingOther
- [ ] Flux.SampleTimeout
- [ ] Flux.SingleOrDefault
- [ ] Flux.SingleOrEmpty

#### Handling errors

- [ ] Timeout
- [x] OnErrorReturn
- [ ] OnErrorResume
- [ ] OnErrorMap
- [ ] Retry
- [ ] RetryWhen
- [ ] Flux.OnBackpressureError
- [ ] Flux.OnBackpressureBuffer
- [ ] Flux.OnBackpressureDrop
- [ ] Flux.OnBackpressureLatest

#### Working with time

- [ ] Elapsed
- [ ] Timestamp
- [ ] Timeout
- [ ] Interval
- [ ] Mono.Delay
- [ ] Mono.DelayElement
- [ ] Flux.DelayElements
- [ ] DelaySubscription

#### Splitting a Flux

- [ ] Flux.Window
- [ ] Flux.WindowPeriod
- [ ] Flux.WindowTimeout
- [ ] Flux.WindowUntil
- [ ] Flux.WindowWhile
- [ ] Flux.WindowUsingOther
- [ ] Flux.WindowWhen
- [ ] Flux.Buffer
- [ ] Flux.BufferPeriod
- [ ] Flux.BufferTimeout
- [ ] Flux.BufferUntil
- [ ] Flux.BufferWhile
- [ ] Flux.BufferWhen
- [ ] Flux.BufferUsingOther
- [ ] Flux.GroupBy

#### Synchronizing

- [x] Flux.BlockFirst
- [x] Flux.BlockFirstTimeout
- [x] Flux.BlockLast
- [x] Flux.BlockLastTimeout
- [x] Mono.Block
- [x] Mono.BlockTimeout

### TODO

- Add schedule periodic and schedule after to schedulers and add ability to insert virtual clock ( this will be useful in tests)
- How to split up tests for normal and scalar flux/mono?
- Fix locking for flatMaps
- Move most docs to godoc, except some examples and "how to choose an operator"
- NoneSignal() ?
- Performance benchmarks