package cesium

import (
	"log"
	"os"
	"time"

	"github.com/DusanKasan/cesium/flux"
)

// Subscriber will receive 3, 4 and complete. We need to sleep at the end so
// that the execution doesn't end and we can wait for the emissions.
func Example() {
	var subscriber Subscriber

	flux.
		Just(1, 2, 3).
		Map(func(t T) T {
			return t.(int) + 1
		}).
		Filter(func(t T) bool {
			return t.(int) > 2
		}).
		Subscribe(subscriber)

	time.Sleep(time.Second)
}

// Output the lifecycle of this Flux to the supplied log and block until the
// flux is terminated.
func Example_blocking() {
	flux.
		Just(1, 2, 3).
		Map(func(t T) T {
			return t.(int) + 1
		}).
		Filter(func(t T) bool {
			return t.(int) > 2
		}).
		Log(log.New(os.Stdout, "", 0)).
		BlockLast()
}
