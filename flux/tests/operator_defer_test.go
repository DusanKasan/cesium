package tests

import (
	"testing"
	"time"

	"math"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestDefer(t *testing.T) {
	publisher := flux.Defer(func() cesium.Publisher {
		return flux.FromSlice([]cesium.T{1, 2, 3})
	})

	verifier.
		Create(publisher).
		ExpectNext(1, 2, 3).
		ExpectComplete().
		Verify(t)
}

func TestDeferSubscriptionBehaviour(t *testing.T) {
	publisher := flux.Defer(func() cesium.Publisher {
		return flux.FromSlice([]cesium.T{1, 2, 3})
	})

	verifier.
		Create(publisher).
		ExpectNext(1, 2).
		ThenCancel().
		ThenAwait(time.Millisecond).
		ExpectNextCount(0).
		Verify(t)
}

func TestDeferUnbounded(t *testing.T) {
	publisher := flux.Defer(func() cesium.Publisher {
		return flux.FromSlice([]cesium.T{1, 2, 3})
	})

	verifier.
		Create(publisher).
		ThenRequest(math.MaxInt64).
		ExpectNextCount(3).
		ExpectComplete().
		Verify(t)
}
