package tests

import (
	"testing"

	"errors"

	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestDoOnComplete(t *testing.T) {
	received := false

	publisher := flux.
		Just(1, 2, 3).
		DoOnComplete(func() {
			received = true
		})

	verifier.
		Create(publisher).
		ExpectNext(1).
		Then(func() {
			if received {
				t.Errorf("OnComplete called too soon")
			}
		}).
		ExpectNext(2, 3).
		ExpectComplete().
		Then(func() {
			if !received {
				t.Errorf("OnComplete not called")
			}
		}).
		Verify(t)
}

func TestDoOnCompleteNotCalledWhenNotCompleted(t *testing.T) {
	received := false

	err := errors.New("err")

	publisher := flux.
		Error(err).
		DoOnComplete(func() {
			received = true
		})

	verifier.
		Create(publisher).
		ExpectError(err).
		Then(func() {
			if received {
				t.Errorf("OnComplete called when it should not have been")
			}
		}).
		Verify(t)
}
