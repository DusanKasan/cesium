package tests

import (
	"testing"
	"time"

	"math"

	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestJust(t *testing.T) {
	verifier.
		Create(flux.Just(1, 2, 3)).
		ExpectNext(1, 2, 3).
		ExpectComplete().
		Verify(t)
}

func TestJustSubscriptionBehaviour(t *testing.T) {
	verifier.
		Create(flux.Just(1, 2, 3)).
		ExpectNext(1, 2).
		ThenCancel().
		ThenAwait(time.Millisecond).
		ExpectNextCount(0).
		Verify(t)
}

func TestJustUnbounded(t *testing.T) {
	verifier.
		Create(flux.Just(1, 2, 3)).
		ThenRequest(math.MaxInt64).
		ExpectNextCount(3).
		ExpectComplete().
		Verify(t)
}
