package tests

import (
	"testing"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestAny(t *testing.T) {
	publisher := flux.
		FromSlice([]cesium.T{1, 2, 3}).
		Any(func(t cesium.T) bool {
			return t.(int) > 1
		})

	verifier.Create(publisher).
		ExpectNext(true).
		ExpectComplete().
		Verify(t)

	publisher = flux.
		FromSlice([]cesium.T{1, 2, 3}).
		Any(func(t cesium.T) bool {
			return t.(int) > 3
		})

	verifier.Create(publisher).
		ExpectNext(false).
		ExpectComplete().
		Verify(t)
}

func TestAnyScalarFlux(t *testing.T) {
	publisher := flux.
		Just(1).
		Any(func(t cesium.T) bool {
			return t.(int) > 2
		})

	verifier.
		Create(publisher).
		ExpectNext(false).
		ExpectComplete()

	publisher = flux.
		Just(3).
		Any(func(t cesium.T) bool {
			return t.(int) > 2
		})

	verifier.
		Create(publisher).
		ExpectNext(true).
		ExpectComplete()
}
