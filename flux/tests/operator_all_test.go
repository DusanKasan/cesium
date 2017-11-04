package tests

import (
	"testing"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestAll(t *testing.T) {
	publisher := flux.
		FromSlice([]cesium.T{1, 2, 3}).
		All(func(t cesium.T) bool {
			return t.(int) > 0
		})

	verifier.
		Create(publisher).
		ExpectNext(true).
		ExpectComplete()
}

func TestAllFalse(t *testing.T) {
	publisher := flux.
		FromSlice([]cesium.T{1, 2, 3}).
		All(func(t cesium.T) bool {
			return t.(int) > 2
		})

	verifier.
		Create(publisher).
		ExpectNext(false).
		ExpectComplete()
}
