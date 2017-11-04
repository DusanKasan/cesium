package tests

import (
	"testing"

	"time"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/mono"
	"github.com/DusanKasan/cesium/verifier"
)

func TestFromCallable(t *testing.T) {
	publisher := mono.FromCallable(func() cesium.T {
		return 1
	})

	verifier.
		Create(publisher).
		ExpectNext(1).
		ExpectComplete().
		Verify(t)
}

func TestFromCallableCancellation(t *testing.T) {
	publisher := mono.FromCallable(func() cesium.T {
		return 1
	})

	verifier.
		Create(publisher).
		ThenCancel().
		ThenRequest(1).
		ThenAwait(time.Millisecond).
		ExpectNextCount(0).
		Verify(t)
}
