package tests

import (
	"testing"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestHasElement(t *testing.T) {
	publisher := flux.
		FromSlice([]cesium.T{1, 2, 3}).
		HasElement(2)

	verifier.Create(publisher).
		ExpectNext(true).
		ExpectComplete().
		Verify(t)

	publisher = flux.
		FromSlice([]cesium.T{1, 2, 3}).
		HasElement(4)

	verifier.Create(publisher).
		ExpectNext(false).
		ExpectComplete().
		Verify(t)
}

func TestHasElementScalarFlux(t *testing.T) {
	publisher := flux.
		Just(1).
		HasElement(1)

	verifier.Create(publisher).
		ExpectNext(true).
		ExpectComplete().
		Verify(t)

	publisher = flux.
		Just(1).
		HasElement(4)

	verifier.Create(publisher).
		ExpectNext(false).
		ExpectComplete().
		Verify(t)
}
