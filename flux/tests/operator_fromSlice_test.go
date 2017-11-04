package tests

import (
	"testing"
	"time"

	"math"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestFromSlice(t *testing.T) {
	verifier.
		Create(flux.FromSlice([]cesium.T{1, 2, 3})).
		ExpectNext(1, 2, 3).
		ExpectComplete().
		Verify(t)
}

func TestFromSliceSubscriptionBehaviour(t *testing.T) {
	verifier.
		Create(flux.FromSlice([]cesium.T{1, 2, 3})).
		ExpectNext(1, 2).
		ThenCancel().
		ThenAwait(time.Millisecond).
		ExpectNextCount(0).
		Verify(t)
}

func TestFromSliceUnbounded(t *testing.T) {
	verifier.
		Create(flux.FromSlice([]cesium.T{1, 2, 3})).
		ThenRequest(math.MaxInt64).
		ExpectNextCount(3).
		ExpectComplete().
		Verify(t)
}
