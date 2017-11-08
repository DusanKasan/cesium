package tests

import (
	"testing"

	"math"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestFromChannel(t *testing.T) {
	c := make(chan cesium.T)

	go func() {
		c <- 1
		c <- 2
		c <- 3
		close(c)
	}()

	verifier.
		Create(flux.FromChannel(c)).
		ExpectNext(1, 2, 3).
		ThenRequest(1).
		ExpectComplete().
		Verify(t)
}

func TestFromChannelUnbounded(t *testing.T) {
	c := make(chan cesium.T)

	go func() {
		c <- 1
		c <- 2
		c <- 3
		close(c)
	}()

	verifier.
		Create(flux.FromChannel(c)).
		ThenRequest(math.MaxInt64).
		ExpectNext(1, 2, 3).
		ExpectComplete().
		Verify(t)
}
