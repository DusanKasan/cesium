package tests

import (
	"testing"
	"time"

	"sync"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/mono"
	"github.com/DusanKasan/cesium/verifier"
)

func TestUsing(t *testing.T) {
	var cleanedUp cesium.T
	cleanupMux := sync.Mutex{}

	publisher := mono.Using(
		func() cesium.T { return 1 },
		func(t cesium.T) cesium.Mono { return mono.Just(t) },
		func(t cesium.T) {
			cleanupMux.Lock()
			cleanedUp = t
			cleanupMux.Unlock()
		},
	)
	verifier.
		Create(publisher).
		ExpectNext(1).
		ExpectComplete().
		Then(func() {
			cleanupMux.Lock()
			if cleanedUp != 1 {
				t.Errorf("Cleanup not called or called with wrong value. Expected %v, Got: %v", 1, cleanedUp)
			}
			cleanupMux.Unlock()
		}).
		Verify(t)
}

func TestUsingSubscriptionBehaviour(t *testing.T) {
	var cleanedUp cesium.T
	cleanupMux := sync.Mutex{}

	publisher := mono.Using(
		func() cesium.T { return 1 },
		func(t cesium.T) cesium.Mono { return mono.Just(t) },
		func(t cesium.T) {
			cleanupMux.Lock()
			cleanedUp = t
			cleanupMux.Unlock()
		},
	)

	verifier.
		Create(publisher).
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
