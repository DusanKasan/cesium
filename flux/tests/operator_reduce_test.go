package tests

import (
	"testing"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestReduce(t *testing.T) {
	publisher := flux.
		FromSlice([]cesium.T{1, 2, 3}).
		Reduce(func(a cesium.T, b cesium.T) cesium.T {
			return a.(int) + b.(int)
		})

	verifier.
		Create(publisher).
		ExpectNext(6).
		ExpectComplete().
		Verify(t)
}

func TestReduceScalarFLux(t *testing.T) {
	publisher := flux.
		Just(1).
		Reduce(func(a cesium.T, b cesium.T) cesium.T {
			return a.(int) + b.(int)
		})

	verifier.
		Create(publisher).
		ExpectNext(1).
		ExpectComplete().
		Verify(t)
}
