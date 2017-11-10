package tests

import (
	"testing"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/mono"
	"github.com/DusanKasan/cesium/verifier"
)

func TestFlatMapMany(t *testing.T) {
	f := mono.Just(1).
		FlatMapMany(func(t cesium.T) cesium.Publisher {
			return flux.Just(t, t.(int)+1)
		})

	verifier.
		Create(f).
		ExpectNext(1, 2).
		ExpectComplete().
		Verify(t)
}
