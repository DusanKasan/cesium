package mono_test

import (
	"errors"

	"fmt"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/mono"
)

// The subscriber will receive item 1 and then an onComplete signal
func ExampleJust() {
	var subscriber cesium.Subscriber

	mono.Just(1).Subscribe(subscriber)
}

// The subscriber will receive item 1 and then an onComplete signal
func ExampleJustOrEmpty() {
	var subscriber cesium.Subscriber

	mono.JustOrEmpty(1).Subscribe(subscriber)
}

// The subscriber will receive no item and an onComplete signal
func ExampleJustOrEmpty_WithNil() {
	var subscriber cesium.Subscriber

	mono.JustOrEmpty(nil).Subscribe(subscriber)
}

// The subscriber will receive item 1 and then an onComplete signal
func ExampleFromCallable() {
	var subscriber cesium.Subscriber

	mono.FromCallable(func() cesium.T {
		return 1
	}).Subscribe(subscriber)
}

// The subscriber will receive no item and an onComplete signal
func ExampleEmpty() {
	var subscriber cesium.Subscriber

	mono.Empty().Subscribe(subscriber)
}

// The subscriber will receive no item or termination signal
func ExampleNever() {
	var subscriber cesium.Subscriber

	mono.Never().Subscribe(subscriber)
}

// The subscriber will receive item 4 and an onComplete signal after which
// "Range mono finished, dispose of the resource (3)" is printed.
//
// Internally the return value of the first argument will be used as an input to
// the second, producing a mono of Just(4). The subscriber is subscribed to
// this mono. Then when the subscribe reaches the termination signal, the third
// argument is executed and "Range mono finished, dispose of the resource (3)"
// is printed.
func ExampleUsing() {
	var subscriber cesium.Subscriber

	mono.Using(
		func() cesium.T {
			return 3
		},
		func(t cesium.T) cesium.Mono {
			return mono.Just(t.(int) + 1)
		},
		func(t cesium.T) {
			fmt.Println("Just mono finished, dispose of the resource (3)")
		},
	).Subscribe(subscriber)
}

// The subscriber will receive item 1 and an onComplete signal. If the
// subscribe would not be able to keep up with how fast we are emitting the item
// , the item would be buffered.
func ExampleCreate() {
	var subscriber cesium.Subscriber

	mono.Create(func(sink cesium.MonoSink) {
		sink.CompleteWith(1)
	}).Subscribe(subscriber)
}

// The subscriber will receive item 1 and an onComplete signal.
func ExampleDefer() {
	var subscriber cesium.Subscriber

	mono.Defer(func() cesium.Mono {
		return mono.Just(1)
	}).Subscribe(subscriber)
}

// The subscriber will receive no item and an onError signal containing the
// supplied error.
func ExampleError() {
	var subscriber cesium.Subscriber

	mono.Error(errors.New("error")).Subscribe(subscriber)
}

// The subscriber will receive item 1 and an onComplete signal.
func ExampleFromChannel() {
	c := make(chan cesium.T)
	var subscriber cesium.Subscriber

	mono.FromChannel(c).Subscribe(subscriber)

	c <- 1
	close(c)
}
