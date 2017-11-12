package verifier_test

import (
	"testing"

	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func Example() {
	t := &testing.T{}

	publisher := flux.Just(1, 2, 3)

	verifier.
		Create(publisher).
		ExpectNext(1, 2, 3).
		ExpectComplete().
		Verify(t)
}
