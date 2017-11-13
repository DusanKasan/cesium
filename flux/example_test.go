package flux_test

import (
	"time"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
)

// The subscriber will receive items 1, 2, 3 and then a onComplete signal
func ExampleJust() {
	var subscriber cesium.Subscriber

	flux.Just(1, 2, 3).Subscribe(subscriber)

	time.Sleep(time.Second)
}
