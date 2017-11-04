package tests

import (
	"testing"

	"errors"

	"github.com/DusanKasan/cesium/mono"
	"github.com/DusanKasan/cesium/verifier"
)

func TestDoOnSuccess(t *testing.T) {
	received := false

	publisher := mono.
		Just(1).
		DoOnSuccess(func() {
			received = true
		})

	verifier.
		Create(publisher).
		Then(func() {
			if received {
				t.Errorf("OnSuccess called too soon")
			}
		}).
		ExpectNext(1).
		ExpectComplete().
		Then(func() {
			if !received {
				t.Errorf("OnSuccess not called")
			}
		}).
		Verify(t)
}

func TestDoOnSuccessNotCalledWhenNotSuccessd(t *testing.T) {
	received := false

	err := errors.New("err")

	publisher := mono.
		Error(err).
		DoOnSuccess(func() {
			received = true
		})

	verifier.
		Create(publisher).
		ExpectError(err).
		Then(func() {
			if received {
				t.Errorf("OnSuccess called when it should not have been")
			}
		}).
		Verify(t)
}
