package tests

import (
	"errors"
	"testing"
	"time"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/mono"
	"github.com/DusanKasan/cesium/verifier"
)

func TestDoFinallyOnComplete(t *testing.T) {
	finallyFlag := false
	finallyLock := make(chan bool)

	publisher := mono.Just(1).DoFinally(func() {
		finallyFlag = true
		finallyLock <- true
	})

	verifier.
		Create(publisher).
		ExpectNext(1).
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

	publisher := mono.Error(err).DoFinally(func() {
		finallyFlag = true
		finallyLock <- true
	})

	verifier.
		Create(publisher).
		ThenRequest(1).
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

	publisher := mono.Create(func(s cesium.MonoSink) {
		time.Sleep(time.Second)
		s.Complete()
	}).DoFinally(func() {
		finallyFlag = true
	})

	verifier.
		Create(publisher).
		ThenCancel().
		Then(func() {
			if !finallyFlag {
				t.Errorf("Finally body not executed")
			}
		}).
		Verify(t)
}
