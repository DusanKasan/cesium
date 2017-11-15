package flux_test

import (
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

// The subscriber will receive items int64(1), int64(2), int64(3) and then an onComplete signal
func ExampleRange() {
	var subscriber cesium.Subscriber

	flux.Range(1, 3).Subscribe(subscriber)
}

// The subscriber will receive no items and an onComplete signal
func ExampleEmpty() {
	var subscriber cesium.Subscriber

	flux.Just(1, 2, 3).Subscribe(subscriber)
}

// The subscriber will receive no items or termination signal
func ExampleNever() {
	var subscriber cesium.Subscriber

	flux.Never().Subscribe(subscriber)
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
