package tests

import (
	"testing"
	"time"

	"math"

	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestRange(t *testing.T) {
	verifier.
		Create(flux.Range(2, 3)).
		ExpectNext(int64(2), int64(3), int64(4)).
		ExpectComplete().
		Verify(t)
}

func TestRangeSubscriptionBehaviour(t *testing.T) {
	verifier.
		Create(flux.Range(2, 3)).
		ExpectNext(int64(2), int64(3)).
		ThenCancel().
		ThenRequest(1).
		ThenAwait(time.Millisecond).
		ExpectNextCount(0).
		Verify(t)
}

func TestRangeUnbounded(t *testing.T) {
	verifier.
		Create(flux.Range(2, 3)).
		ThenRequest(math.MaxInt64).
		ExpectNextCount(3).
		ExpectComplete().
		Verify(t)
}
