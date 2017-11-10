package tests

import (
	"testing"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestFlatMap(t *testing.T) {
	c := make(chan cesium.T)
	received := make(chan bool)

	f := flux.FromChannel(c).
		FlatMap(func(a cesium.T) cesium.Publisher {
			return flux.Just(a, a.(int)+1)
		})

	go func() {
		<-received
		c <- 1
		<-received
		c <- 2
		close(c)
	}()

	received <- true

	verifier.
		Create(f).
		ExpectNext(1, 2).
		Then(func() {
			received <- true
		}).
		ExpectNext(2, 3).
		ExpectComplete().
		Verify(t)
}
