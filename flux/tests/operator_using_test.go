package tests

import (
	"testing"
	"time"

	"sync"

	"math"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestUsing(t *testing.T) {
	var cleanedUp cesium.T
	in := make(chan bool)
	out := make(chan bool)

	publisher := flux.Using(
		func() cesium.T { return 1 },
		func(t cesium.T) cesium.Publisher { return flux.Just(t, t.(int)+1) },
		func(t cesium.T) {
			<-in
			cleanedUp = t
			out <- true
		},
	)

	verifier.
		Create(publisher).
		ExpectNext(1, 2).
		ExpectComplete().
		Then(func() {
			// We need to give Using time to cleanup
			in <- true
			<-out

			if cleanedUp != 1 {
				t.Errorf("Cleanup not called or called with wrong value. Expected %v, Got: %v", 1, cleanedUp)
			}
		}).
		Verify(t)

}

func TestUsingSubscriptionBehaviour(t *testing.T) {
	var cleanedUp cesium.T
	cleanupMux := sync.Mutex{}

	publisher := flux.Using(
		func() cesium.T { return 1 },
		func(t cesium.T) cesium.Publisher { return flux.Just(t, t.(int)+1) },
		func(t cesium.T) {
			cleanupMux.Lock()
			cleanedUp = t
			cleanupMux.Unlock()
		},
	)

	verifier.
		Create(publisher).
		ExpectNext(1).
		Then(func() {
			if cleanedUp != nil {
				t.Errorf("Cleanup called too soon")
			}
		}).
		ThenCancel().
		Then(func() {
			if cleanedUp != 1 {
				t.Errorf("Cleanup not called or called with wrong value. Expected %v, Got: %v", 1, cleanedUp)
			}
		}).
		ThenRequest(1).
		ThenAwait(time.Millisecond).
		ExpectNextCount(0).
		Verify(t)
}

func TestUsingUnbounded(t *testing.T) {
	var cleanedUp cesium.T
	in := make(chan bool)
	out := make(chan bool)

	publisher := flux.Using(
		func() cesium.T { return 1 },
		func(t cesium.T) cesium.Publisher { return flux.Just(t, t.(int)+1) },
		func(t cesium.T) {
			<-in
			cleanedUp = t
			out <- true
		},
	)

	verifier.
		Create(publisher).
		ThenRequest(math.MaxInt64).
		ExpectNextCount(2).
		ExpectComplete().
		Then(func() {
			// We need to give Using time to cleanup
			in <- true
			<-out

			if cleanedUp != 1 {
				t.Errorf("Cleanup not called or called with wrong value. Expected %v, Got: %v", 1, cleanedUp)
			}
		}).
		Verify(t)
}
