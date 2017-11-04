package tests

import (
	"testing"

	"errors"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestDoFinallyOnComplete(t *testing.T) {
	finallyFlag := false
	finallyLock := make(chan bool)

	publisher := flux.
		FromSlice([]cesium.T{1, 2, 3}).
		DoFinally(func() {
			finallyFlag = true
			finallyLock <- true
		})

	verifier.
		Create(publisher).
		ExpectNext(1, 2, 3).
		ExpectComplete().
		Then(func() {
			<-finallyLock

			if !finallyFlag {
				t.Errorf("Finally body not executed")
			}
		}).
		Verify(t)

}

func TestDoFinallyOnError(t *testing.T) {
	finallyFlag := false
	finallyLock := make(chan bool)
	err := errors.New("err")

	publisher := flux.
		Error(err).
		DoFinally(func() {
			finallyFlag = true
			finallyLock <- true
		})

	verifier.
		Create(publisher).
		ExpectError(err).
		Then(func() {
			<-finallyLock

			if !finallyFlag {
				t.Errorf("Finally body not executed")
			}
		}).
		Verify(t)
}

func TestDoFinallyOnCancel(t *testing.T) {
	finallyFlag := false

	publisher := flux.Just(1, 2).DoFinally(func() {
		finallyFlag = true
	})

	verifier.
		Create(publisher).
		ExpectNext(1).
		ThenCancel().
		Then(func() {
			if !finallyFlag {
				t.Errorf("Finally body not executed")
			}
		}).
		Verify(t)
}
