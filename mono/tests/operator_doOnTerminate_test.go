package tests

import (
	"errors"
	"testing"

	"github.com/DusanKasan/cesium/mono"
	"github.com/DusanKasan/cesium/verifier"
)

func TestDoOnTerminateWhenError(t *testing.T) {
	err := errors.New("err")
	received := false

	publisher := mono.
		Error(err).
		DoOnTerminate(func() {
			received = true
		})

	verifier.
		Create(publisher).
		ExpectError(err).
		Then(func() {
			if !received {
				t.Errorf("On error not called")
			}
		}).
		Verify(t)
}

func TestDoOnTerminateWhenComplete(t *testing.T) {
	received := false

	publisher := mono.
		Just(1).
		DoOnTerminate(func() {
			received = true
		})

	verifier.
		Create(publisher).
		Then(func() {
			if received {
				t.Errorf("On error called too soon")
			}
		}).
		ExpectNext(1).
		ExpectComplete().
		Then(func() {
			if !received {
				t.Errorf("On error not called")
			}
		}).
		Verify(t)
}
