package tests

import (
	"testing"
	"time"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/mono"
	"github.com/DusanKasan/cesium/verifier"
)

func TestDefer(t *testing.T) {
	publisher := mono.Defer(func() cesium.Mono {
		return mono.Just(1)
	})

	verifier.
		Create(publisher).
		ExpectNext(1).
		ExpectComplete().
		Verify(t)
}

func TestDeferSubscriptionBehaviour(t *testing.T) {
	publisher := mono.Defer(func() cesium.Mono {
		return mono.Just(1)
	})

	verifier.
		Create(publisher).
		ThenCancel().
		ThenRequest(1).
		ThenAwait(time.Millisecond).
		ExpectNextCount(0).
		Verify(t)
}
