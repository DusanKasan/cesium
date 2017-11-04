package verifier

import (
	"sync"
	"testing"

	"time"

	"math"

	"github.com/DusanKasan/cesium"
)

type expectation struct {
	expectationType string
	value           cesium.T
	err             error
	callback        func(cesium.T)
}

type subscription struct {
	onCancel           func()
	onRequest          func(int64)
	onRequestUnbounded func()
}

func (s *subscription) Cancel() {
	s.onCancel()
}

func (s *subscription) Request(n int64) {
	s.onRequest(n)
}

func (s *subscription) RequestUnbounded() {
	s.onRequestUnbounded()
}

type bufferObserver struct {
	buffer       []cesium.MaterializedEmission
	subscription cesium.Subscription
	mux          sync.Mutex
}

func (o *bufferObserver) OnNext(t cesium.T) {
	o.mux.Lock()
	o.buffer = append(o.buffer, cesium.MaterializedEmission{
		EventType: "next",
		Value:     t,
	})
	o.mux.Unlock()
}

func (o *bufferObserver) OnComplete() {
	o.mux.Lock()
	o.buffer = append(o.buffer, cesium.MaterializedEmission{
		EventType: "complete",
	})
	o.mux.Unlock()
}

func (o *bufferObserver) OnError(err error) {
	o.mux.Lock()
	o.buffer = append(o.buffer, cesium.MaterializedEmission{
		EventType: "error",
		Err:       err,
	})
	o.mux.Unlock()
}

func (o *bufferObserver) OnSubscribe(subscription cesium.Subscription) {
	o.mux.Lock()
	o.subscription = subscription
	o.mux.Unlock()
}

func (o *bufferObserver) GetSubscription() cesium.Subscription {
	o.mux.Lock()
	s := o.subscription
	o.mux.Unlock()

	return s
}

func (o *bufferObserver) Pop() *cesium.MaterializedEmission {
	var emission *cesium.MaterializedEmission

	o.mux.Lock()
	if len(o.buffer) > 0 {
		emission = &o.buffer[0]
		o.buffer = o.buffer[1:]
	}
	o.mux.Unlock()
	return emission
}

func (o *bufferObserver) BufferedCount() int {
	o.mux.Lock()
	buffered := len(o.buffer)
	o.mux.Unlock()
	return buffered
}

func (o *bufferObserver) BufferedNextCount() int {
	o.mux.Lock()
	size := len(o.buffer)
	if size > 0 {
		if o.buffer[size-1].EventType != "next" {
			size = size - 1
		}
	}
	o.mux.Unlock()
	return size
}

func (o *bufferObserver) PurgeBuffer() {
	o.mux.Lock()
	o.buffer = []cesium.MaterializedEmission{}
	o.mux.Unlock()
}

type stepVerifier struct {
	publisher    cesium.Publisher
	expectations []expectation
	timeout      time.Duration
}

func Create(p cesium.Publisher) *stepVerifier {
	return &stepVerifier{publisher: p, timeout: time.Millisecond * 200}
}

func (sv *stepVerifier) AndTimeout(duration time.Duration) *stepVerifier {
	sv.timeout = duration

	return sv
}

func (sv *stepVerifier) ExpectNext(items ...cesium.T) *stepVerifier {
	for _, i := range items {
		sv.expectations = append(sv.expectations, expectation{
			expectationType: "next",
			value:           i,
		})
	}

	return sv
}

func (sv *stepVerifier) ExpectComplete() *stepVerifier {
	sv.expectations = append(sv.expectations, expectation{
		expectationType: "complete",
	})

	return sv
}

func (sv *stepVerifier) ExpectError(err error) *stepVerifier {
	sv.expectations = append(sv.expectations, expectation{
		expectationType: "error",
		err:             err,
	})

	return sv
}

func (sv *stepVerifier) ThenCancel() *stepVerifier {
	sv.expectations = append(sv.expectations, expectation{
		expectationType: "cancel",
	})

	return sv
}

func (sv *stepVerifier) ThenAwait(duration time.Duration) *stepVerifier {
	sv.expectations = append(sv.expectations, expectation{
		expectationType: "await",
		value:           duration,
	})

	return sv
}

func (sv *stepVerifier) ThenRequest(i int64) *stepVerifier {
	sv.expectations = append(sv.expectations, expectation{
		expectationType: "request",
		value:           i,
	})

	return sv
}

func (sv *stepVerifier) ExpectNextCount(i int) *stepVerifier {
	sv.expectations = append(sv.expectations, expectation{
		expectationType: "nextCount",
		value:           i,
	})

	return sv
}

func (sv *stepVerifier) Then(f func()) *stepVerifier {
	sv.expectations = append(sv.expectations, expectation{
		expectationType: "then",
		value:           f,
	})

	return sv
}

func (sv *stepVerifier) WithTimeout(duration time.Duration) *stepVerifier {
	sv.timeout = duration
	return sv
}

func (sv *stepVerifier) Verify(t *testing.T) {
	observer := &bufferObserver{}
	s := sv.publisher.Subscribe(observer)
	pendingRequests := int64(0)

	obsSub := observer.GetSubscription()
	if s != obsSub {
		t.Errorf("The returned subscription is not the same the Subscriber received. Returned: %#v, Subscriber got: %#v", s, obsSub)
		return
	}

	mux := sync.Mutex{}
	canceled := false

	subs := &subscription{
		onRequestUnbounded: func() {
			s.RequestUnbounded()
		},
		onRequest: func(i int64) {
			s.Request(i)
		},
		onCancel: func() {
			mux.Lock()
			canceled = true
			s.Cancel()
			mux.Unlock()
		},
	}

	previousNextBufferedCount := 0
	for i, e := range sv.expectations {
		switch e.expectationType {
		case "cancel":
			subs.Cancel()
			continue
		case "await":
			time.Sleep(e.value.(time.Duration))
			continue
		case "request":
			subs.Request(e.value.(int64))
			if e.value.(int64) == math.MaxInt64 {
				pendingRequests = math.MaxInt64 //TODO: There should be a bool holding isUnboundedMode on/off
			} else {
				pendingRequests = pendingRequests + e.value.(int64)
			}
			continue
		case "then":
			e.value.(func())()
			continue
		case "nextCount":
			start := time.Now()
			for {
				actualSize := observer.BufferedNextCount() - previousNextBufferedCount

				if actualSize == e.value.(int) {
					break
				}

				if start.Add(sv.timeout).Before(time.Now()) {
					t.Errorf("Wrong number of emissions from last expectation. Expected: %v, Got: %v", e.value.(int), actualSize)
					return
				}
			}

			for i := e.value.(int); i > 0; i-- {
				observer.Pop()
			}

			pendingRequests = pendingRequests - int64(e.value.(int))
			previousNextBufferedCount = observer.BufferedNextCount()

			continue
		case "next":
			if pendingRequests == 0 {
				pendingRequests = 1
				subs.Request(1)
			}
		}

		var emission *cesium.MaterializedEmission

		for emission == nil {
			mux.Lock()
			if canceled {
				if i < len(sv.expectations)-1 {
					t.Errorf("Subscription cancelled before it finished.")
					return
					mux.Unlock()
				}
				mux.Unlock()
			}
			mux.Unlock()
			// return with timeout error after a period
			emission = observer.Pop()
		}

		if emission.EventType == "next" {
			pendingRequests = pendingRequests - 1
		}

		mux.Lock()
		if canceled {
			if i < len(sv.expectations)-1 {
				t.Errorf("Subscription cancelled before it finished.")
				return
				mux.Unlock()
			}
			mux.Unlock()
		}
		mux.Unlock()
		if !compareExpectationAndEmission(e, *emission, t) {
			subs.Cancel()
			return
		}

		previousNextBufferedCount = observer.BufferedNextCount()
	}
}

func compareExpectationAndEmission(expected expectation, actual cesium.MaterializedEmission, t *testing.T) bool {
	switch expected.expectationType {
	case "next":
		if actual.EventType != "next" {
			t.Errorf("Invalid emission type! Expected: %v, Got: %v", "next", actual.EventType)
			return false
		}

		if actual.Value != expected.value {
			t.Errorf("Invalid value! Expected: %v, Got: %v", expected.value, actual.Value)
			return false
		}
	case "complete":
		if actual.EventType != "complete" {
			t.Errorf("Invalid emission type! Expected: %v, Got: %v", "complete", actual.EventType)
			return false
		}
	case "error":
		if actual.EventType != "error" {
			t.Errorf("Invalid emission type! Expected: %v, Got: %v", "error", actual.EventType)
			return false
		}

		if actual.Err != expected.err {
			t.Errorf("Invalid error! Expected: %v, Got: %v", expected.err, actual.Err)
			return false
		}
	}

	return true
}
