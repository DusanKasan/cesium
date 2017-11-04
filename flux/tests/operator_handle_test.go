package tests

import (
	"testing"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestHandle(t *testing.T) {
	f := flux.Just(1, 2, 3).
		Handle(func(t cesium.T, sink cesium.SynchronousSink) {
			switch t {
			case 1:
				sink.Next(10)
			case 2:
				sink.Next(15)
			default:
				sink.Complete()
			}
		})

	verifier.
		Create(f).
		ExpectNext(10, 15).
		ThenRequest(1).
		ExpectComplete().
		Verify(t)
}

func TestHandleWhenNoEmissions(t *testing.T) {
	f := flux.Just(1, 2, 3).
		Handle(func(t cesium.T, sink cesium.SynchronousSink) {
			switch t {
			case 1:
				sink.Next(10)
			case 2:
				sink.Next(15)
			}
		})

	verifier.
		Create(f).
		ExpectNext(10, 15).
		ThenRequest(1).
		ExpectError(cesium.NoEmissionOnSynchronousSinkError).
		Verify(t)
}
