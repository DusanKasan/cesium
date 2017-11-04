package tests

import (
	"errors"
	"testing"

	"github.com/DusanKasan/cesium/mono"
	"github.com/DusanKasan/cesium/verifier"
)

func TestDoOnError(t *testing.T) {
	err := errors.New("err")
	received := false

	publisher := mono.
		Error(err).
		DoOnError(func(e error) {
			if err == e {
				received = true
			}
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

func TestDoOnErrorNotCalledWhenNoError(t *testing.T) {
	err := errors.New("err")
	received := false

	publisher := mono.
		Just(1).
		DoOnError(func(e error) {
			if err == e {
				received = true
			}
		})

	verifier.
		Create(publisher).
		ExpectNext(1).
		ExpectComplete().
		Then(func() {
			if received {
				t.Errorf("On error called when it should not have been")
			}
		}).
		Verify(t)
}
