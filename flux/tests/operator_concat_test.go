package tests

import (
	"testing"

	"errors"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestConcat(t *testing.T) {
	publisher := flux.
		FromSlice([]cesium.T{1, 2}).
		Concat(flux.FromSlice([]cesium.T{
			flux.FromSlice([]cesium.T{3, 4}),
			flux.FromSlice([]cesium.T{5, 6}),
		}))

	verifier.
		Create(publisher).
		ExpectNext(1, 2, 3, 4, 5, 6).
		ExpectComplete().
		Verify(t)
}

func TestConcatOnErrorInSubPublisher(t *testing.T) {
	err := errors.New("err")

	publisher := flux.
		FromSlice([]cesium.T{1, 2}).
		Concat(flux.FromSlice([]cesium.T{
			flux.Error(err),
			flux.FromSlice([]cesium.T{5, 6}),
		}))

	verifier.
		Create(publisher).
		ExpectNext(1, 2).
		ExpectError(err).
		Verify(t)
}
