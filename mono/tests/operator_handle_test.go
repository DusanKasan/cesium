package tests

import (
	"testing"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/mono"
	"github.com/DusanKasan/cesium/verifier"
)

func TestHandle(t *testing.T) {
	f := mono.Just(1).
		Handle(func(t cesium.T, sink cesium.SynchronousSink) {
			sink.Next(10)
		})

	verifier.
		Create(f).
		ExpectNext(10).
		ExpectComplete().
		Verify(t)
}

func TestHandleWhenNoEmissions(t *testing.T) {
	f := mono.Just(1).
		Handle(func(t cesium.T, sink cesium.SynchronousSink) {
		})

	verifier.
		Create(f).
		ThenRequest(1).
		ExpectError(cesium.NoEmissionOnSynchronousSinkError).
		Verify(t)
}

func TestHandleScalarFlux(t *testing.T) {
	f := mono.Just(1).
		Handle(func(t cesium.T, sink cesium.SynchronousSink) {
			sink.Complete()
		})

	verifier.
		Create(f).
		ThenRequest(1).
		ExpectComplete().
		Verify(t)
}
