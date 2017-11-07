package tests

import (
	"testing"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

// Tests filter on just(), proving that it works with ConditionalSubscriber
func TestFilterSyncSourceWithMicroOperatorFusion(t *testing.T) {
	f := flux.Just(2, 30, 22, 5, 60, 1).
		Filter(func(a cesium.T) bool {
			return a.(int) > 10
		})

	verifier.
		Create(f).
		ExpectNext(30, 22, 60).
		ThenRequest(1).
		ExpectComplete().
		Verify(t)
}

func TestFilterAsyncSource(t *testing.T) {
	f := flux.Create(
		func(sink cesium.FluxSink) {
			sink.Next(2)
			sink.Next(30)
			sink.Next(22)
			sink.Next(5)
			sink.Next(60)
			sink.Next(1)
			sink.Complete()
		},
		flux.OverflowStrategyBuffer).
		Filter(func(a cesium.T) bool {
			return a.(int) > 10
		})

	verifier.
		Create(f).
		ExpectNext(30, 22, 60).
		ThenRequest(1).
		ExpectComplete().
		Verify(t)
}
