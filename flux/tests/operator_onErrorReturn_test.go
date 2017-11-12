package tests

import (
	"testing"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
	"github.com/pkg/errors"
)

func TestOnErrorReturn(t *testing.T) {
	publisher := flux.
		Create(func(sink cesium.FluxSink) {
			sink.Next(1)
			sink.Error(errors.New("err"))
		}, flux.OverflowStrategyBuffer).
		OnErrorReturn(2)

	verifier.Create(publisher).
		ExpectNext(1, 2).
		ExpectComplete().
		Verify(t)
}
