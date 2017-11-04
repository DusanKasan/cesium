package tests

import (
	"testing"

	"github.com/DusanKasan/cesium/mono"
	"github.com/DusanKasan/cesium/verifier"
)

func TestDoOnRequest(t *testing.T) {
	requested := int64(0)

	publisher := mono.
		Just(1).
		DoOnRequest(func(n int64) {
			requested = requested + n
		})

	verifier.
		Create(publisher).
		Then(func() {
			if requested != 0 {
				t.Errorf("OnRequested called incorectly: Expected: %v, Got: %v", 0, requested)
			}
		}).
		ExpectNext(1).
		Then(func() {
			if requested != 1 {
				t.Errorf("OnRequested called incorectly: Expected: %v, Got: %v", 1, requested)
			}
		}).
		ExpectComplete().
		Verify(t)

}
