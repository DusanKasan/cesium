package mono

import (
	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/internal"
)

// Just creates new cesium.Mono that emits the supplied item and completes.
func Just(t cesium.T) cesium.Mono {
	return internal.MonoJust(t)
}

// JustOrEmpty creates new cesium.Mono that emits the supplied item if it's non-nil
// and completes, otherwise just completes.
func JustOrEmpty(t cesium.T) cesium.Mono {
	return internal.MonoJustOrEmpty(t)
}

// FromCallable creates new cesium.Mono that emits the item returned from the supplied function. If the function
// returns nil, the returned Mono completes empty.
func FromCallable(f func() cesium.T) cesium.Mono {
	return internal.MonoFromCallable(f)
}

// Empty creates new cesium.Mono that emits no items and completes normally.
func Empty() cesium.Mono {
	return internal.MonoEmpty()
}

// Empty creates new cesium.Mono that emits no items and completes with error.
func Error(err error) cesium.Mono {
	return internal.MonoError(err)
}

// Never creates new cesium.Mono that emits no items and never completes.
func Never() cesium.Mono {
	return internal.MonoNever()
}

// Defer creates new cesium.Mono by subscribing to the Mono returned from the supplied
// factory function for each subscribtion.
func Defer(f func() cesium.Mono) cesium.Mono {
	return internal.MonoDefer(f)
}

// Using Uses a resource, generated by a supplier for each individual Subscriber,
// while streaming the value from a Mono derived from the same resource and makes
// sure the resource is released if the sequence terminates or the Subscriber cancels.
func Using(resourceSupplier func() cesium.T, sourceSupplier func(cesium.T) cesium.Mono, resourceCleanup func(cesium.T)) cesium.Mono {
	return internal.MonoUsing(resourceSupplier, sourceSupplier, resourceCleanup)
}

// Create allows you to programmatically create a cesium.Mono with the
// capability of emitting multiple elements in a synchronous or asynchronous
// manner through the mono.Sink API
func Create(f func(cesium.MonoSink)) cesium.Mono {
	return internal.MonoCreate(f)
}