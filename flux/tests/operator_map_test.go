package tests

import (
	"testing"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestMap(t *testing.T) {
	publisher := flux.
		FromSlice([]cesium.T{1, 2, 3}).
		Map(func(a cesium.T) cesium.T {
			return 10 * a.(int)
		})

	verifier.
		Create(publisher).
		ExpectNext(10, 20, 30).
		ExpectComplete().
		Verify(t)
}
