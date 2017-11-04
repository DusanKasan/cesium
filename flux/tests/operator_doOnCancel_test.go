package tests

import (
	"testing"

	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestDoOnCancel(t *testing.T) {
	received := false

	publisher := flux.
		Just(1, 2, 3).
		DoOnCancel(func() {
			received = true
		})

	verifier.
		Create(publisher).
		ExpectNext(1).
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

	publisher := flux.
		Just(1, 2).
		DoOnCancel(func() {
			received = true
		})

	verifier.
		Create(publisher).
		ExpectNext(1, 2).
		ExpectComplete().
		Then(func() {
			if received {
				t.Errorf("OnCancel called when it should not have been")
			}
		}).
		Verify(t)
}
