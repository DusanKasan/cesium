package internal

import "github.com/DusanKasan/cesium"

func ErrorSignal(err error) cesium.Signal {
	return &signal{
		signalType: cesium.SignalTypeOnError,
		err:        err,
	}
}

func NextSignal(t cesium.T) cesium.Signal {
	return &signal{
		signalType: cesium.SignalTypeOnNext,
		item:       t,
	}
}

func CompleteSignal() cesium.Signal {
	return &signal{
		signalType: cesium.SignalTypeOnComplete,
	}
}

func SubscribtionSignal(s cesium.Subscription) cesium.Signal {
	return &signal{
		signalType:   cesium.SignalTypeOnSubscribe,
		subscription: s,
	}
}

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
