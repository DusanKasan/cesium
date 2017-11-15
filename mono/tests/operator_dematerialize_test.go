package tests

import (
	"testing"

	"errors"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/mono"
	"github.com/DusanKasan/cesium/verifier"
)

type signal struct {
	signalType   cesium.SignalType
	err          error
	item         cesium.T
	subscription cesium.Subscription
}

func (s *signal) Accept(subscriber cesium.Subscriber) {
	switch s.signalType {
	case cesium.SignalTypeOnSubscribe:
		subscriber.OnSubscribe(s.subscription)
	case cesium.SignalTypeOnNext:
		subscriber.OnNext(s.item)
	case cesium.SignalTypeOnComplete:
		subscriber.OnComplete()
	case cesium.SignalTypeOnError:
		subscriber.OnError(s.err)
	}
}

func (s *signal) Item() cesium.T {
	return s.item
}

func (s *signal) Subscription() cesium.Subscription {
	return s.subscription
}

func (s *signal) Error() error {
	return s.err
}

func (s *signal) Type() cesium.SignalType {
	return s.signalType
}

func (s *signal) IsOnSubscribe() bool {
	return s.signalType == cesium.SignalTypeOnSubscribe
}

func (s *signal) IsOnNext() bool {
	return s.signalType == cesium.SignalTypeOnNext
}

func (s *signal) IsOnComplete() bool {
	return s.signalType == cesium.SignalTypeOnComplete
}

func (s *signal) IsOnError() bool {
	return s.signalType == cesium.SignalTypeOnError
}

func (s *signal) IsTerminal() bool {
	return s.IsOnComplete() || s.IsOnError()
}

func errorSignal(err error) cesium.Signal {
	return &signal{
		signalType: cesium.SignalTypeOnError,
		err:        err,
	}
}

func nextSignal(t cesium.T) cesium.Signal {
	return &signal{
		signalType: cesium.SignalTypeOnNext,
		item:       t,
	}
}

func completeSignal() cesium.Signal {
	return &signal{
		signalType: cesium.SignalTypeOnComplete,
	}
}

func TestDematerialize(t *testing.T) {
	publisher := mono.
		Just(nextSignal(1)).
		Dematerialize()

	verifier.
		Create(publisher).
		ExpectNext(1).
		ThenRequest(1).
		ExpectComplete().
		Verify(t)

	err := errors.New("err")
	publisher = mono.
		Just(errorSignal(err)).
		Dematerialize()

	verifier.
		Create(publisher).
		ThenRequest(1).
		ExpectError(err).
		Verify(t)

}
