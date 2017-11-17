package tests

import (
	"testing"

	"errors"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestOnErrorMap(t *testing.T) {
	err := errors.New("err")
	replacementError := errors.New("repl")

	f := flux.
		Create(
			func(sink cesium.FluxSink) {
				sink.Next(1)
				sink.Error(err)
			},
			flux.OverflowStrategyBuffer).
		OnErrorMap(
			func(e error) error {
				return replacementError
			})

	verifier.
		Create(f).
		ExpectNext(1).
		ExpectError(replacementError).
		Verify(t)
}
