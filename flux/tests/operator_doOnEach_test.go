package tests

import (
	"testing"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
	"github.com/pkg/errors"
)

func TestDoOnEach(t *testing.T) {

	var buffer []cesium.Signal

	publisher := flux.
		Just(1).
		DoOnEach(func(s cesium.Signal) {
			buffer = append(buffer, s)
		})

	verifier.
		Create(publisher).
		ExpectNext(1).
		ThenRequest(1).
		ExpectComplete().
		Then(func() {
			if len(buffer) != 3 {
				t.Errorf("Unexpected number of signales received. Expected: %v, Got: %v", 3, len(buffer))
				return
			}

			if buffer[0].Type() != cesium.SignalTypeOnSubscribe {
				t.Errorf("Unexpected signal type at index 0 received. Expected: %v, Got: %v", cesium.SignalTypeOnSubscribe, buffer[0].Type())
				return
			}

			if buffer[1].Type() != cesium.SignalTypeOnNext {
				t.Errorf("Unexpected signal type at index 1 received. Expected: %v, Got: %v", cesium.SignalTypeOnNext, buffer[0].Type())
				return
			}

			if buffer[2].Type() != cesium.SignalTypeOnComplete {
				t.Errorf("Unexpected signal type at index 2 received. Expected: %v, Got: %v", cesium.SignalTypeOnComplete, buffer[0].Type())
				return
			}
		}).
		Verify(t)
}

func TestDoOnEachWhenError(t *testing.T) {

	var buffer []cesium.Signal

	err := errors.New("err")
	publisher := flux.
		Error(err).
		DoOnEach(func(s cesium.Signal) {
			buffer = append(buffer, s)
		})

	verifier.
		Create(publisher).
		ExpectError(err).
		Then(func() {

			if len(buffer) != 2 {
				t.Errorf("Unexpected number of signales received. Expected: %v, Got: %v", 2, len(buffer))
				return
			}

			if buffer[0].Type() != cesium.SignalTypeOnSubscribe {
				t.Errorf("Unexpected signal type at index 0 received. Expected: %v, Got: %v", cesium.SignalTypeOnSubscribe, buffer[0].Type())
				return
			}

			if buffer[1].Type() != cesium.SignalTypeOnError {
				t.Errorf("Unexpected signal type at index 0 received. Expected: %v, Got: %v", cesium.SignalTypeOnError, buffer[0].Type())
				return
			}
		}).
		Verify(t)
}
