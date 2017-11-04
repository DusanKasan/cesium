package tests

import (
	"testing"

	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestDoOnRequest(t *testing.T) {
	requested := int64(0)

	publisher := flux.
		Just(1, 2, 3).
		DoOnRequest(func(n int64) {
			requested = requested + n
		})

	verifier.
		Create(publisher).
		ExpectNext(1).
		Then(func() {
			if requested != 1 {
				t.Errorf("OnRequested called incorectly: Expected: %v, Got: %v", 1, requested)
			}
		}).
		ExpectNext(2, 3).
		ExpectComplete().
		Then(func() {
			if requested != 3 {
				t.Errorf("OnRequested called incorectly: Expected: %v, Got: %v", 3, requested)
			}
		}).
		Verify(t)

}
