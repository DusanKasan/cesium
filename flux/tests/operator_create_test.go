package tests

import (
	"errors"
	"sync"
	"testing"
	"time"

	"math"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestCreate(t *testing.T) {
	disposed := false
	mux := sync.Mutex{}

	publisher := flux.Create(func(s cesium.FluxSink) {
		s.OnDispose(func() {
			mux.Lock()
			disposed = true
			mux.Unlock()
		})

		s.Next(1)
		s.Next(2)
		s.Next(3)
		s.Complete()
	}, flux.OverflowStrategyBuffer)

	verifier.
		Create(publisher).
		ExpectNext(1, 2, 3).
		ExpectComplete().
		Then(func() {
			mux.Lock()
			if !disposed {
				t.Errorf("Dispose not signalled to Sink")
			}
			mux.Unlock()
		}).
		Verify(t)

}

func TestCreateSubscriptionBehaviour(t *testing.T) {
	cancelled := false
	disposed := false
	requested := int64(0)
	mux := sync.Mutex{}

	publisher := flux.Create(func(s cesium.FluxSink) {
		s.OnCancel(func() {
			mux.Lock()
			cancelled = true
			mux.Unlock()
		})

		s.OnDispose(func() {
			mux.Lock()
			disposed = true
			mux.Unlock()
		})

		s.OnRequest(func(n int64) {
			mux.Lock()
			requested = requested + n
			mux.Unlock()
		})

		s.Next(1)
		s.Next(2)
		s.Next(2)
		s.Next(3)
		s.Complete()
	}, flux.OverflowStrategyBuffer)

	verifier.
		Create(publisher).
		ExpectNext(1, 2).
		ThenCancel().
		ThenRequest(1).
		ThenAwait(time.Millisecond).
		ExpectNextCount(0).
		Then(func() {
			mux.Lock()
			if !cancelled {
				t.Errorf("Cancellation not signalled to Sink")
			}

			if !disposed {
				t.Errorf("Disposed not signalled to Sink")
			}

			if requested != 2 {
				t.Errorf("RequestFunc signalled incorectly, Expected: %v, Got: %v", 2, requested)
			}
			mux.Unlock()
		}).
		Verify(t)
}

func TestCreateDisposedOnError(t *testing.T) {
	mux := sync.Mutex{}
	disposed := false
	err := errors.New("err")

	publisher := flux.Create(func(s cesium.FluxSink) {
		s.OnDispose(func() {
			mux.Lock()
			disposed = true
			mux.Unlock()
		})

		s.Next(1)
		s.Next(2)
		s.Next(3)
		s.Error(err)
	}, flux.OverflowStrategyBuffer)

	verifier.
		Create(publisher).
		ExpectNext(1, 2, 3).
		ExpectError(err).
		Then(func() {
			mux.Lock()
			if !disposed {
				t.Errorf("Disposed not signalled to Sink")
			}
			mux.Unlock()
		}).
		Verify(t)
}

func TestCreateWithDrop(t *testing.T) {
	in := make(chan bool)
	out := make(chan bool)

	publisher := flux.Create(func(s cesium.FluxSink) {
		s.Next(1)
		s.Next(2)

		in <- true
		<-out

		s.Next(3)
		s.Complete()
	}, flux.OverflowStrategyDrop)

	verifier.
		Create(publisher).
		Then(func() {
			<-in
		}).
		ThenRequest(1).
		Then(func() { out <- true }).
		ExpectNext(3).
		ExpectComplete().
		Verify(t)
}

func TestCreateWithError(t *testing.T) {
	in := make(chan bool)

	publisher := flux.Create(func(s cesium.FluxSink) {
		<-in
		s.Next(1)
		s.Next(2)
		s.Next(3)
		s.Complete()
	}, flux.OverflowStrategyError)

	verifier.
		Create(publisher).
		ThenRequest(2).
		Then(func() { in <- true }).
		ExpectNextCount(2).
		ExpectError(cesium.DownstreamUnableToKeepUpError).
		Verify(t)
}

func TestCreateWithIgnore(t *testing.T) {
	publisher := flux.Create(func(s cesium.FluxSink) {
		s.Next(1)
		s.Next(2)
		s.Next(3)
		s.Complete()
	}, flux.OverflowStrategyIgnore)

	verifier.
		Create(publisher).
		ExpectNextCount(3).
		ExpectComplete().
		Verify(t)
}

func TestCreateWithBufferUnbounded(t *testing.T) {
	publisher := flux.Create(func(s cesium.FluxSink) {
		s.Next(1)
		s.Next(2)
		s.Next(3)
		s.Complete()
	}, flux.OverflowStrategyBuffer)

	verifier.
		Create(publisher).
		ThenRequest(math.MaxInt64).
		ExpectNextCount(3).
		ExpectComplete().
		Verify(t)
}

func TestCreateWithErrorUnbounded(t *testing.T) {
	in := make(chan bool)

	publisher := flux.Create(func(s cesium.FluxSink) {
		<-in
		s.Next(1)
		s.Next(2)
		s.Next(3)
		s.Complete()
	}, flux.OverflowStrategyError)

	verifier.
		Create(publisher).
		ThenRequest(math.MaxInt64).
		Then(func() { in <- true }).
		ExpectNextCount(3).
		ExpectComplete().
		Verify(t)
}

func TestCreateWithDropUnbounded(t *testing.T) {
	out := make(chan bool)

	publisher := flux.Create(func(s cesium.FluxSink) {
		<-out

		s.Next(1)
		s.Next(2)
		s.Next(3)
		s.Complete()
	}, flux.OverflowStrategyDrop)

	verifier.
		Create(publisher).
		ThenRequest(math.MaxInt64).
		Then(func() { out <- true }).
		ExpectNextCount(3).
		ExpectComplete().
		Verify(t)
}
