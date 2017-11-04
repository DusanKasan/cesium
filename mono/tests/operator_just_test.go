package tests

import (
	"testing"

	"time"

	"github.com/DusanKasan/cesium/mono"
	"github.com/DusanKasan/cesium/verifier"
)

func TestJust(t *testing.T) {
	verifier.
		Create(mono.Just(1)).
		ExpectNext(1).
		ExpectComplete().
		Verify(t)
}

func TestJustCancellation(t *testing.T) {
	verifier.
		Create(mono.Just(1)).
		ThenCancel().
		ThenRequest(1).
		ThenAwait(time.Millisecond).
		ExpectNextCount(0).
		Verify(t)
}
