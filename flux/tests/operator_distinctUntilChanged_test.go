package tests

import (
	"testing"

	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestDistinctUntilChanged(t *testing.T) {
	f := flux.Just(1, 2, 2, 3, 1, 2, 3, 3).DistinctUntilChanged()

	verifier.
		Create(f).
		ExpectNext(1, 2, 3, 1, 2, 3).
		ThenRequest(1).
		ExpectComplete().
		Verify(t)
}

func TestDistinctUntilChangedScalarFlux(t *testing.T) {
	f := flux.Just(1).DistinctUntilChanged()

	verifier.
		Create(f).
		ExpectNext(1).
		ExpectComplete().
		Verify(t)
}
