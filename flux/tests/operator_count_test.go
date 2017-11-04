package tests

import (
	"testing"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestCount(t *testing.T) {
	publisher := flux.
		FromSlice([]cesium.T{1, 2, 3}).
		Count()

	verifier.
		Create(publisher).
		ExpectNext(int64(3)).
		ExpectComplete().
		Verify(t)

}
