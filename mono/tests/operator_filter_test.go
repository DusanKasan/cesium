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
}

func TestFilterWithEmptyResult(t *testing.T) {
	publisher := mono.
		Just(3).
		Filter(func(a cesium.T) bool {
			return a.(int) > 4
		})

	verifier.
		Create(publisher).
		ThenRequest(1).
		ExpectComplete().
		Verify(t)
}
