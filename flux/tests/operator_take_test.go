package tests

import (
	"testing"

	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestTake(t *testing.T) {
	f := flux.Just(1, 2, 3, 4, 5).Take(3)

	verifier.
		Create(f).
		ExpectNext(1, 2, 3).
		ExpectComplete().
		Verify(t)
}
