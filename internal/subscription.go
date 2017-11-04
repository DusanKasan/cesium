package internal

import (
	"math"
	"sync"

	"github.com/DusanKasan/cesium"
)

type Subscription struct {
	CancelFunc  func()
	RequestFunc func(int64)
}

func (s *Subscription) Cancel() {
	s.CancelFunc()
}

func (s *Subscription) Request(n int64) {
	s.RequestFunc(n)
}

func (s *Subscription) RequestUnbounded() {
	s.RequestFunc(math.MaxInt64)
}

type BufferedProxySubscription struct {
	mux          sync.Mutex
	subscription cesium.Subscription

	cancelled bool
	requested int64
}

func (s *BufferedProxySubscription) SetSubscription(subscription cesium.Subscription) {
	s.mux.Lock()
	s.subscription = subscription
	if s.cancelled {
		s.subscription.Cancel()
	}

	if s.requested > 0 {
		s.subscription.Request(s.requested)
	}
	s.mux.Unlock()
}

func (s *BufferedProxySubscription) Cancel() {
	s.mux.Lock()
	if s.subscription != nil {
		s.subscription.Cancel()
	} else {
		s.cancelled = true
	}
	s.mux.Unlock()
}

func (s *BufferedProxySubscription) Request(n int64) {
	s.mux.Lock()
	if s.subscription != nil {
		s.subscription.Request(n)
	} else {
		s.requested = s.requested + n
	}
	s.mux.Unlock()
}
