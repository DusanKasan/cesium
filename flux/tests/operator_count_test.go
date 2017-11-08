package tests

import (
	"testing"

	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestCount(t *testing.T) {
	publisher := flux.
		Just(1, 2, 3).
		Count()

	verifier.
		Create(publisher).
		ExpectNext(int64(3)).
		ExpectComplete().
		Verify(t)
}

func TestCountScalarFlux(t *testing.T) {
	publisher := flux.
		Just(1).
		Count()

	verifier.
		Create(publisher).
		ExpectNext(int64(1)).
		ExpectComplete().
		Verify(t)
}
