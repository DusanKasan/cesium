package tests

import (
	"testing"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestScan(t *testing.T) {
	publisher := flux.
		FromSlice([]cesium.T{1, 2, 3}).
		Scan(func(a cesium.T, b cesium.T) cesium.T {
			return a.(int) + b.(int)
		})

	verifier.
		Create(publisher).
		ExpectNext(1, 3, 6).
		ExpectComplete().
		Verify(t)
}

func TestScanScalarFlux(t *testing.T) {
	publisher := flux.
		Just(1).
		Scan(func(a cesium.T, b cesium.T) cesium.T {
			return a.(int) + b.(int)
		})

	verifier.
		Create(publisher).
		ExpectNext(1).
		ExpectComplete().
		Verify(t)
}
