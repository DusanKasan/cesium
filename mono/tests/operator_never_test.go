package tests

import (
	"testing"
	"time"

	"github.com/DusanKasan/cesium/mono"
	"github.com/DusanKasan/cesium/verifier"
)

func TestNever(t *testing.T) {
	verifier.
		Create(mono.Never()).
		ThenRequest(1).
		ThenAwait(time.Millisecond).
		ExpectNextCount(0).
		Verify(t)
}
