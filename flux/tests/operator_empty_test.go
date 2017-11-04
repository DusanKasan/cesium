package tests

import (
	"testing"

	"time"

	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestEmpty(t *testing.T) {
	verifier.
		Create(flux.Empty()).
		ThenRequest(1).
		ThenAwait(time.Millisecond).
		ExpectNextCount(0).
		Verify(t)
}
