package tests

import (
	"testing"

	"time"

	"github.com/DusanKasan/cesium/mono"
	"github.com/DusanKasan/cesium/verifier"
)

func TestJustOrEmpty(t *testing.T) {
	verifier.
		Create(mono.JustOrEmpty(1)).
		ExpectNext(1).
		ExpectComplete().
		Verify(t)
}

func TestJustOrEmptyWithNil(t *testing.T) {
	verifier.
		Create(mono.JustOrEmpty(nil)).
		ExpectComplete().
		Verify(t)
}

func TestJustOrEmptyCancellation(t *testing.T) {
	verifier.
		Create(mono.JustOrEmpty(1)).
		ThenCancel().
		ThenRequest(1).
		ThenAwait(time.Millisecond).
		ExpectNextCount(0).
		Verify(t)
}
