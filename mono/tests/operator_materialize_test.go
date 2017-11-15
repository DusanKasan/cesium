package tests

import (
	"testing"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/mono"
	"github.com/DusanKasan/cesium/verifier"
)

func TestMaterialize(t *testing.T) {
	publisher := mono.
		Just(1).
		Materialize()

	verifier.
		Create(publisher).
		ExpectNextMatches(func(i cesium.T) bool {
			return i.(cesium.Signal).Item() == 1
		}).
		ExpectNextMatches(func(i cesium.T) bool {
			return i.(cesium.Signal).IsOnComplete()
		}).
		ThenRequest(1).
		ExpectComplete().
		Verify(t)
}
