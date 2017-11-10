package tests

import (
	"testing"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/mono"
	"github.com/DusanKasan/cesium/verifier"
)

func TestFlatMap(t *testing.T) {
	f := mono.Just(1).
		FlatMap(func(t cesium.T) cesium.Mono {
			return mono.Just(t.(int) + 1)
		})

	verifier.
		Create(f).
		ExpectNext(2).
		ExpectComplete().
		Verify(t)
}
