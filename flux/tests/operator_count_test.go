package tests

import (
	"testing"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestCount(t *testing.T) {
	publisher := flux.
		FromSlice([]cesium.T{1, 2, 3}).
		Count()

	verifier.
		Create(publisher).
		ExpectNext(int64(3)).
		ExpectComplete().
		Verify(t)
}

// We know that Just with single item will resolve to a scalar flux, and that
// Count on this flux supports macro fusion and will not fire the whole Rx
// infrastructure. We are testing the correctness here.
func TestCountOperatorMacroFusionOnScalarFlux(t *testing.T) {
	publisher := flux.
		Just(1).
		Count()

	verifier.
		Create(publisher).
		ExpectNext(int64(1)).
		ExpectComplete().
		Verify(t)
}
