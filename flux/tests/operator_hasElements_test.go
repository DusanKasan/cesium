package tests

import (
	"testing"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestHasElements(t *testing.T) {
	publisher := flux.
		FromSlice([]cesium.T{1, 2, 3}).
		HasElements()

	verifier.
		Create(publisher).
		ExpectNext(true).
		ExpectComplete().
		Verify(t)

	publisher = flux.
		Empty().
		HasElements()

	verifier.
		Create(publisher).
		ExpectNext(false).
		ExpectComplete().
		Verify(t)
}

func TestHasElementsScalarFlux(t *testing.T) {
	publisher := flux.
		Just(1).
		HasElements()

	verifier.
		Create(publisher).
		ExpectNext(true).
		ExpectComplete().
		Verify(t)
}
