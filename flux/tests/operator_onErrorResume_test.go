package tests

import (
	"testing"

	"errors"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestOnErrorResume(t *testing.T) {
	err := errors.New("err")

	f := flux.
		Create(
			func(sink cesium.FluxSink) {
				sink.Next(1)
				sink.Error(err)
			},
			flux.OverflowStrategyBuffer).
		OnErrorResume(
			func(e error) bool {
				return e == err
			},
			flux.Just(2, 3))

	verifier.
		Create(f).
		ExpectNext(1, 2, 3).
		ExpectComplete().
		Verify(t)
}
