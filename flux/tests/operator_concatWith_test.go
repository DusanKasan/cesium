package tests

import (
	"testing"

	"errors"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestConcatWith(t *testing.T) {
	publisher := flux.
		Just(1, 2).
		ConcatWith(
			flux.FromSlice([]cesium.T{3, 4}),
			flux.FromSlice([]cesium.T{5, 6}),
		)

	verifier.
		Create(publisher).
		ExpectNext(1, 2, 3, 4, 5, 6).
		ExpectComplete().
		Verify(t)
}

func TestConcatWithOnErrorInSubPublisher(t *testing.T) {
	err := errors.New("err")

	publisher := flux.
		Just(1, 2).
		ConcatWith(
			flux.Error(err),
			flux.FromSlice([]cesium.T{3, 4}),
		)

	verifier.
		Create(publisher).
		ExpectNext(1, 2).
		ExpectError(err).
		//ExpectNothing(time.Duration)
		Verify(t)
}
