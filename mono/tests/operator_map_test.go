package tests

import (
	"testing"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/mono"
	"github.com/DusanKasan/cesium/verifier"
)

func TestMap(t *testing.T) {
	publisher := mono.
		Just(1).
		Map(func(a cesium.T) cesium.T {
			return 10 * a.(int)
		})

	verifier.
		Create(publisher).
		ExpectNext(10).
		ExpectComplete().
		Verify(t)
}
