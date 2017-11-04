package tests

import (
	"testing"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/mono"
	"github.com/DusanKasan/cesium/verifier"
)

func TestDoOnSubscribe(t *testing.T) {
	received := false

	publisher := mono.
		Just(1).
		DoOnSubscribe(func(s cesium.Subscription) {
			received = true
		})

	verifier.
		Create(publisher).
		Then(func() {
			if !received {
				t.Errorf("OnSubscribe not called")
			}
		}).
		ExpectNext(1).
		ExpectComplete().
		Verify(t)

}
