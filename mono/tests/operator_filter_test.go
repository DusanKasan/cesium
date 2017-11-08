package tests

import (
	"testing"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/mono"
	"github.com/DusanKasan/cesium/verifier"
)

func TestFilter(t *testing.T) {
	publisher := mono.
		Just(5).
		Filter(func(a cesium.T) bool {
			return a.(int) > 4
		})

	verifier.
		Create(publisher).
		ExpectNext(5).
		ExpectComplete().
		Verify(t)

	publisher2 := mono.
		Just(3).
		Filter(func(a cesium.T) bool {
			return a.(int) > 4
		})

	verifier.
		Create(publisher2).
		ThenRequest(1).
		ExpectComplete().
		Verify(t)
}

func TestFilterMacroFusionScalarMono(t *testing.T) {
	f := mono.Just(11).
		Filter(func(a cesium.T) bool {
			return a.(int) > 10
		})

	verifier.
		Create(f).
		ExpectNext(11).
		ThenRequest(1).
		ExpectComplete().
		Verify(t)

	g := mono.Just(1).
		Filter(func(a cesium.T) bool {
			return a.(int) > 10
		})

	verifier.
		Create(g).
		ThenRequest(1).
		ExpectComplete().
		Verify(t)
}
