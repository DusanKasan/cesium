package tests

import (
	"testing"
	"time"

	"sync"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/mono"
	"github.com/DusanKasan/cesium/verifier"
	"github.com/pkg/errors"
)

func TestCreate(t *testing.T) {
	disposed := false
	mux := sync.Mutex{}

	publisher := mono.Create(func(s cesium.MonoSink) {
		s.OnDispose(func() {
			mux.Lock()
			disposed = true
			mux.Unlock()
		})

		s.CompleteWith(1)
	})

	verifier.
		Create(publisher).
		ExpectNext(1).
		ExpectComplete().
		Verify(t)

	mux.Lock()
	if !disposed {
		t.Errorf("Disposed not signalled to Sink")
	}
	mux.Unlock()
}

func TestCreateCancellationBehaviour(t *testing.T) {
	cancelled := false
	disposed := false
	mux := sync.Mutex{}

	publisher := mono.Create(func(s cesium.MonoSink) {
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

		s.CompleteWith(1)
	})

	verifier.
		Create(publisher).
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
			mux.Unlock()
		}).
		Verify(t)
}

func TestCreateOnError(t *testing.T) {
	disposed := false
	mux := sync.Mutex{}

	err := errors.New("err")

	publisher := mono.Create(func(s cesium.MonoSink) {
		s.OnDispose(func() {
			mux.Lock()
			disposed = true
			mux.Unlock()
		})

		s.Error(err)
	})

	verifier.
		Create(publisher).
		ThenRequest(1).
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
