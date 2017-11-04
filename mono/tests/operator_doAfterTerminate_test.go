package tests

import (
	"errors"
	"testing"
	"time"

	"github.com/DusanKasan/cesium/mono"
	"github.com/DusanKasan/cesium/verifier"
)

func TestDoAfterTerminateWhenError(t *testing.T) {
	err := errors.New("err")
	receive := make(chan bool)

	publisher := mono.
		Error(err).
		DoAfterTerminate(func() {
			receive <- true
		})

	verifier.
		Create(publisher).
		ExpectError(err).
		Then(func() {
			timeout := time.After(time.Second)
			for {
				select {
				case <-timeout:
					t.Errorf("On error not called")
					return
				case <-receive:
					return
				}
			}
		}).
		Verify(t)
}

func TestDoAfterTerminateWhenComplete(t *testing.T) {
	receive := make(chan bool)

	publisher := mono.
		Just(1).
		DoAfterTerminate(func() {
			receive <- true
		})

	verifier.
		Create(publisher).
		Then(func() {
			select {
			case <-receive:
				t.Errorf("On error called too soon")
			default:
			}
		}).
		ExpectNext(1).
		ExpectComplete().
		Then(func() {
			timeout := time.After(time.Second)
			for {
				select {
				case <-timeout:
					t.Errorf("On error not called")
					return
				case <-receive:
					return
				}
			}
		}).
		Verify(t)
}
