package tests

import (
	"testing"

	"errors"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/mono"
	"github.com/DusanKasan/cesium/verifier"
)

func TestConcatWith(t *testing.T) {
	publisher := mono.
		Just(1).
		ConcatWith(
			flux.FromSlice([]cesium.T{2, 3}),
			flux.FromSlice([]cesium.T{4, 5}),
		)

	verifier.
		Create(publisher).
		ExpectNext(1, 2, 3, 4, 5).
		ExpectComplete().
		Verify(t)
}

func TestConcatWithOnErrorInSubPublisher(t *testing.T) {
	err := errors.New("err")

	publisher := mono.
		Just(1).
		ConcatWith(
			mono.Error(err),
			flux.FromSlice([]cesium.T{2, 3}),
		)

	verifier.
		Create(publisher).
		ExpectNext(1).
		ExpectError(err).
		//ExpectNothing(time.Duration)
		Verify(t)
}
