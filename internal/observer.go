package internal

import "github.com/DusanKasan/cesium"

type subscriber struct {
	onNext     func(cesium.T)
	onComplete func()
	onError    func(error)
}

func (s *subscriber) OnNext(t cesium.T) {
	s.onNext(t)
}

func (s *subscriber) OnComplete() {
	s.onComplete()
}

func (s *subscriber) OnError(err error) {
	s.onError(err)
}

func (s *subscriber) OnSubscribe(subscription cesium.Subscription) {
}

func DoObserver(onNext func(cesium.T), onComplete func(), onError func(error)) cesium.Subscriber {
	return &subscriber{
		onNext: func(t cesium.T) {
			onNext(t)
		},
		onComplete: func() {
			onComplete()
		},
		onError: func(err error) {
			onError(err)
		},
	}
}
