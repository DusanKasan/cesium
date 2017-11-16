package flux_test

import (
	"errors"

	"fmt"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
)

// The subscriber will receive items 1, 2, 3 and then an onComplete signal
func ExampleJust() {
	var subscriber cesium.Subscriber

	flux.Just(1, 2, 3).Subscribe(subscriber)
}

// The subscriber will receive items 1, 2, 3 and then an onComplete signal
func ExampleFromSlice() {
	var subscriber cesium.Subscriber

	flux.Just([]cesium.T{1, 2, 3}).Subscribe(subscriber)
}

// The subscriber will receive items int64(1), int64(2), int64(3) and then an
// onComplete signal
func ExampleRange() {
	var subscriber cesium.Subscriber

	flux.Range(1, 3).Subscribe(subscriber)
}

// The subscriber will receive no items and an onComplete signal
func ExampleEmpty() {
	var subscriber cesium.Subscriber

	flux.Empty().Subscribe(subscriber)
}

// The subscriber will receive no items or termination signal
func ExampleNever() {
	var subscriber cesium.Subscriber

	flux.Never().Subscribe(subscriber)
}

// The subscriber will receive items 1, 2, 3 and an onComplete signal after which
// "Range flux finished, dispose of the resource (3)" is printed.
//
// Internally the return value of the first argument will be used as an input to
// the second, producing a flux of Range(1, 3). The subscriber is subscribed to
// this flux. Then when the subscribe reaches the termination signal, the third
// argument is executed and "Range flux finished, dispose of the resource (3)"
// is printed.
func ExampleUsing() {
	var subscriber cesium.Subscriber

	flux.Using(
		func() cesium.T {
			return 3
		},
		func(t cesium.T) cesium.Publisher {
			return flux.Range(1, t.(int))
		},
		func(t cesium.T) {
			fmt.Println("Range flux finished, dispose of the resource (3)")
		},
	).Subscribe(subscriber)
}

// The subscriber will receive items 1, 2, 3 and an onComplete signal. If the
// subscribe would not be able to keep up with how fast we are emitting items
// they will be buffered thanks to the flux.OverflowStrategyBuffer.
func ExampleCreate() {
	var subscriber cesium.Subscriber

	flux.Create(func(sink cesium.FluxSink) {
		sink.Next(1)
		sink.Next(2)
		sink.Next(3)
		sink.Complete()
	}, flux.OverflowStrategyBuffer).Subscribe(subscriber)
}

// The subscriber will receive items 1, 2, 3 and an onComplete signal. You can
// also see that in the last pass when we signal completion the second signal is
// ignored.
func ExampleGenerate() {
	var subscriber cesium.Subscriber

	i := 1
	flux.Generate(func(sink cesium.SynchronousSink) {
		if i > 3 {
			sink.Complete()
		}

		sink.Next(i)
		i++
	}).Subscribe(subscriber)
}

// The subscriber will receive items 1, 2, 3 and an onComplete signal.
func ExampleDefer() {
	var subscriber cesium.Subscriber

	flux.Defer(func() cesium.Publisher {
		return flux.Just(1, 2, 3)
	}).Subscribe(subscriber)
}

// The subscriber will receive no items and an onError signal containing the
// supplied error.
func ExampleError() {
	var subscriber cesium.Subscriber

	flux.Error(errors.New("error")).Subscribe(subscriber)
}

// The subscriber will receive items 1, 2, 3 and an onComplete signal.
func ExampleFromChannel() {
	c := make(chan cesium.T)
	var subscriber cesium.Subscriber

	flux.FromChannel(c).Subscribe(subscriber)

	c <- 1
	c <- 2
	c <- 3
	close(c)
}
