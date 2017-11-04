package tests

import (
	"testing"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestFilter(t *testing.T) {
	f := flux.FromSlice([]cesium.T{2, 30, 22, 5, 60, 1}).
		Filter(func(a cesium.T) bool {
			return a.(int) > 10
		})

	verifier.
		Create(f).
		ExpectNext(30, 22, 60).
		ThenRequest(1).
		ExpectComplete().
		Verify(t)
}
