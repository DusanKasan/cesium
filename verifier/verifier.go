// Package verifier contains the StepVerifier type which provides a declarative
// way of creating a verifiable script for an async Publisher sequence, by
// expressing expectations about the events that will happen upon subscription.
// The verification is triggered by calling the verify(*testing.T) method.
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

type emissionType string

// Represent an onNext emission type in MaterializedEmission.
const EventTypeNext emissionType = "next"

// Represent an onComplete emission type in MaterializedEmission.
const EventTypeComplete emissionType = "complete"

// Represent an onComplete emission type in MaterializedEmission.
const EventTypeError emissionType = "error"

// MaterializedEmission is used in (De)Materialize operator.
type MaterializedEmission struct {
	// The type of emission, EventTypeNext, EventTypeComplete, EventTypeError
	EventType emissionType

	// If the emission was of type "error" the error will be in this field
	Err error

	// If the emission was of type "next" the item will be in this field.
	Value cesium.T
}

type bufferObserver struct {
	buffer       []MaterializedEmission
	subscription cesium.Subscription
	mux          sync.Mutex
}

func (o *bufferObserver) OnNext(t cesium.T) {
	o.mux.Lock()
	o.buffer = append(o.buffer, MaterializedEmission{
		EventType: "next",
		Value:     t,
	})
	o.mux.Unlock()
}

func (o *bufferObserver) OnComplete() {
	o.mux.Lock()
	o.buffer = append(o.buffer, MaterializedEmission{
		EventType: "complete",
	})
	o.mux.Unlock()
}

func (o *bufferObserver) OnError(err error) {
	o.mux.Lock()
	o.buffer = append(o.buffer, MaterializedEmission{
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

func (o *bufferObserver) Pop() *MaterializedEmission {
	var emission *MaterializedEmission

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
	o.buffer = []MaterializedEmission{}
	o.mux.Unlock()
}

type StepVerifier struct {
	publisher    cesium.Publisher
	expectations []expectation
	timeout      time.Duration
}

// Creates a step verifier around the passed Publisher.
func Create(p cesium.Publisher) *StepVerifier {
	return &StepVerifier{publisher: p, timeout: time.Millisecond * 200}
}

// Specifies a timeout for the expectations to come.
func (sv *StepVerifier) AndTimeout(duration time.Duration) *StepVerifier {
	sv.timeout = duration

	return sv
}

// Expect the specified items to be emitted. Handles the requesting internally.
func (sv *StepVerifier) ExpectNext(items ...cesium.T) *StepVerifier {
	for _, i := range items {
		sv.expectations = append(sv.expectations, expectation{
			expectationType: "next",
			value:           i,
		})
	}

	return sv
}

// Expect complete signal to be emitted.
func (sv *StepVerifier) ExpectComplete() *StepVerifier {
	sv.expectations = append(sv.expectations, expectation{
		expectationType: "complete",
	})

	return sv
}

// Expect error signal with the specified error to be emitted.
func (sv *StepVerifier) ExpectError(err error) *StepVerifier {
	sv.expectations = append(sv.expectations, expectation{
		expectationType: "error",
		err:             err,
	})

	return sv
}

// Cancel the underlying subscription.
func (sv *StepVerifier) ThenCancel() *StepVerifier {
	sv.expectations = append(sv.expectations, expectation{
		expectationType: "cancel",
	})

	return sv
}

// Wait for the specified duration before executing next expectation.
func (sv *StepVerifier) ThenAwait(duration time.Duration) *StepVerifier {
	sv.expectations = append(sv.expectations, expectation{
		expectationType: "await",
		value:           duration,
	})

	return sv
}

// Request the specified amount of items from the underlying subscription.
func (sv *StepVerifier) ThenRequest(i int64) *StepVerifier {
	sv.expectations = append(sv.expectations, expectation{
		expectationType: "request",
		value:           i,
	})

	return sv
}

// Checks the specified amount of items to be emitted.
func (sv *StepVerifier) ExpectNextCount(i int) *StepVerifier {
	sv.expectations = append(sv.expectations, expectation{
		expectationType: "nextCount",
		value:           i,
	})

	return sv
}

// Expects the callback to return true for the next item.
func (sv *StepVerifier) ExpectNextMatches(f func(cesium.T) bool) *StepVerifier {
	sv.expectations = append(sv.expectations, expectation{
		expectationType: "nextMatches",
		value:           f,
	})

	return sv
}

// Execute the passed function.
func (sv *StepVerifier) Then(f func()) *StepVerifier {
	sv.expectations = append(sv.expectations, expectation{
		expectationType: "then",
		value:           f,
	})

	return sv
}

// Subscribe to the underlying publisher and start executing the expectation
// chain. Output the errors to the passed T.
func (sv *StepVerifier) Verify(t *testing.T) {
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
		case "next", "nextMatches":
			if pendingRequests == 0 {
				pendingRequests = 1
				subs.Request(1)
			}
		}

		var emission *MaterializedEmission

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

func compareExpectationAndEmission(expected expectation, actual MaterializedEmission, t *testing.T) bool {
	switch expected.expectationType {
	case "nextMatches":
		if actual.EventType != "next" {
			t.Errorf("Invalid emission type! Expected: %v, Got: %v", "next", actual.EventType)
			return false
		}

		if !expected.value.(func(cesium.T) bool)(actual.Value) {
			t.Errorf("ExpectNextMatches returned false for %v", actual.Value)
			return false
		}
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
