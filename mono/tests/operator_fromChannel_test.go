package tests

import (
	"testing"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/mono"
	"github.com/DusanKasan/cesium/verifier"
)

func TestFromChannel(t *testing.T) {
	c := make(chan cesium.T)

	go func() {
		c <- 1
		close(c)
	}()

	verifier.
		Create(mono.FromChannel(c)).
		ExpectNext(1).
		ExpectComplete().
		Verify(t)
}

func TestFromChannelDropsItemsAfterFirst(t *testing.T) {
	c := make(chan cesium.T)

	go func() {
		c <- 1
		c <- 2
		c <- 3
		close(c)
	}()

	verifier.
		Create(mono.FromChannel(c)).
		ExpectNext(1).
		ExpectComplete().
		Verify(t)
}
