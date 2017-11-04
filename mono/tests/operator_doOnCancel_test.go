package tests

import (
	"testing"

	"github.com/DusanKasan/cesium/mono"
	"github.com/DusanKasan/cesium/verifier"
)

func TestDoOnCancel(t *testing.T) {
	received := false

	publisher := mono.
		Just(1).
		DoOnCancel(func() {
			received = true
		})

	verifier.
		Create(publisher).
		ThenCancel().
		Then(func() {
			if !received {
				t.Errorf("OnCancel not called")
			}
		}).
		Verify(t)
}

func TestDoOnCancelNotCalledWhenNotCancelled(t *testing.T) {
	received := false

	publisher := mono.
		Just(1).
		DoOnCancel(func() {
			received = true
		})

	verifier.
		Create(publisher).
		ExpectNext(1).
		ExpectComplete().
		Then(func() {
			if received {
				t.Errorf("OnCancel called when it should not have been")
			}
		}).
		Verify(t)
}
