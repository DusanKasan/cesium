package internal

import (
	"sync"

	"math"

	"github.com/DusanKasan/cesium"
)

type FluxSink struct {
	cancelerMux sync.Mutex
	canceller   cesium.Canceller
	onCancel    func()

	onDisposeMux sync.Mutex
	onDispose    func()

	onRequestMux sync.Mutex
	onRequest    func(int64)

	requestedFromDownstream func() int64
	next                    func(cesium.T)
	complete                func()
	error                   func(error)
	internalRequest         func(int64)

	recordedNumberOfRequests int64
	recordedDispose          bool
}

func (f *FluxSink) Next(t cesium.T) {
	f.next(t)
}

func (f *FluxSink) Complete() {
	f.complete()

	f.onDisposeMux.Lock()
	if f.onDispose != nil {
		f.onDispose()
	} else {
		f.recordedDispose = true
	}
	f.onDisposeMux.Unlock()
}

func (f *FluxSink) Error(err error) {
	f.error(err)

	f.onDisposeMux.Lock()
	if f.onDispose != nil {
		f.onDispose()
	} else {
		f.recordedDispose = true
	}
	f.onDisposeMux.Unlock()
}

func (f *FluxSink) IsCancelled() bool {
	f.cancelerMux.Lock()
	r := f.canceller.IsCancelled()
	f.cancelerMux.Unlock()

	return r
}

func (f *FluxSink) OnCancel(fn func()) {
	onCancel := func() {
		f.Dispose()
		fn()
	}

	f.cancelerMux.Lock()
	if f.canceller.IsCancelled() {
		onCancel()
	}
	f.onCancel = onCancel
	f.cancelerMux.Unlock()
}

func (f *FluxSink) Cancel() {
	f.cancelerMux.Lock()
	if f.onCancel != nil {
		f.onCancel()
	}
	f.cancelerMux.Unlock()

}

func (f *FluxSink) OnDispose(fn func()) {
	f.onDisposeMux.Lock()
	f.onDispose = fn
	if f.recordedDispose {
		f.onDispose()
	}
	f.onDisposeMux.Unlock()
}

func (f *FluxSink) Dispose() {
	f.onDisposeMux.Lock()
	f.onDispose()
	f.onDisposeMux.Unlock()
}

func (f *FluxSink) OnRequest(fn func(n int64)) {
	f.onRequestMux.Lock()
	f.onRequest = fn
	if f.recordedNumberOfRequests > 0 {
		f.onRequest(f.recordedNumberOfRequests)
	}
	f.onRequestMux.Unlock()
}

func (f *FluxSink) Request(n int64) {
	f.cancelerMux.Lock()
	if f.canceller.IsCancelled() {
		f.cancelerMux.Unlock()
		return
	}
	f.cancelerMux.Unlock()

	f.onRequestMux.Lock()
	if f.onRequest != nil {
		f.onRequest(n)
	}
	f.recordedNumberOfRequests = f.recordedNumberOfRequests + n
	f.onRequestMux.Unlock()

	f.internalRequest(n)
}

func (f *FluxSink) RequestedFromDownstream() int64 {
	return f.requestedFromDownstream()
}

func BufferFluxSink(s cesium.Subscriber, c cesium.Canceller) *FluxSink {
	scheduler := SeparateGoroutineScheduler()

	bufferMux := sync.Mutex{}
	buffer := []cesium.MaterializedEmission{}

	requestedMux := sync.Mutex{}
	requested := int64(0)
	unbounded := false

	fs := &FluxSink{
		next: func(t cesium.T) {
			if c.IsCancelled() {
				return
			}

			requestedMux.Lock()
			if unbounded {
				requestedMux.Unlock()
				scheduler.Schedule(func(c cesium.Canceller) {
					s.OnNext(t)
				})
				return
			}

			if requested > 0 {
				requested = requested - 1
				requestedMux.Unlock()
				scheduler.Schedule(func(c cesium.Canceller) {
					s.OnNext(t)
				})
				return
			}
			requestedMux.Unlock()

			bufferMux.Lock()
			buffer = append(buffer, cesium.MaterializedEmission{EventType: "next", Value: t})
			bufferMux.Unlock()
		},
		complete: func() {
			if c.IsCancelled() {
				return
			}

			requestedMux.Lock()
			if unbounded {
				requestedMux.Unlock()
				scheduler.Schedule(func(c cesium.Canceller) {
					s.OnComplete()
				})
				return
			}

			if requested > 0 {
				requested = requested - 1
				requestedMux.Unlock()
				scheduler.Schedule(func(c cesium.Canceller) {
					s.OnComplete()
				})
				return
			}
			requestedMux.Unlock()

			bufferMux.Lock()
			if len(buffer) == 0 {
				scheduler.Schedule(func(c cesium.Canceller) {
					s.OnComplete()
				})
			} else {
				buffer = append(buffer, cesium.MaterializedEmission{EventType: "complete"})
			}
			bufferMux.Unlock()
		},
		error: func(err error) {
			if c.IsCancelled() {
				return
			}

			requestedMux.Lock()
			if unbounded {
				requestedMux.Unlock()
				scheduler.Schedule(func(c cesium.Canceller) {
					s.OnError(err)
				})
				return
			}

			if requested > 0 {
				requested = requested - 1
				requestedMux.Unlock()
				scheduler.Schedule(func(c cesium.Canceller) {
					s.OnError(err)
				})
				return
			}
			requestedMux.Unlock()

			bufferMux.Lock()
			if len(buffer) == 0 {
				scheduler.Schedule(func(c cesium.Canceller) {
					s.OnError(err)
				})
			} else {
				buffer = append(buffer, cesium.MaterializedEmission{EventType: "error", Err: err})
			}
			bufferMux.Unlock()
		},
		internalRequest: func(n int64) {
			if c.IsCancelled() {
				return
			}

			requestedMux.Lock()
			if unbounded {
				return
			}

			if n == math.MaxInt64 {
				unbounded = true
			}

			if unbounded {
				bufferMux.Lock()
				for len(buffer) > 0 {
					emission := buffer[0]
					buffer = buffer[1:]

					switch emission.EventType {
					case "next":
						scheduler.Schedule(func(c cesium.Canceller) {
							s.OnNext(emission.Value)
						})
					case "complete":
						scheduler.Schedule(func(c cesium.Canceller) {
							s.OnComplete()
						})
					case "error":
						scheduler.Schedule(func(c cesium.Canceller) {
							s.OnError(emission.Err)
						})
					}
				}
				bufferMux.Unlock()
				requestedMux.Unlock()
				return
			}

			requested = requested + n
			bufferMux.Lock()
			for len(buffer) > 0 && requested > 0 {
				emission := buffer[0]
				buffer = buffer[1:]
				requested = requested - 1

				switch emission.EventType {
				case "next":
					scheduler.Schedule(func(c cesium.Canceller) {
						s.OnNext(emission.Value)
					})
				case "complete":
					scheduler.Schedule(func(c cesium.Canceller) {
						s.OnComplete()
					})
				case "error":
					scheduler.Schedule(func(c cesium.Canceller) {
						s.OnError(emission.Err)
					})
				}
			}

			if len(buffer) == 1 {
				emission := buffer[0]
				switch emission.EventType {
				case "complete":
					scheduler.Schedule(func(c cesium.Canceller) {
						s.OnComplete()
					})
					buffer = []cesium.MaterializedEmission{}
				case "error":
					scheduler.Schedule(func(c cesium.Canceller) {
						s.OnError(emission.Err)
					})
					buffer = []cesium.MaterializedEmission{}
				}
			}
			bufferMux.Unlock()
			requestedMux.Unlock()
		},
		canceller: c,
	}

	c.OnCancel(fs.Cancel)

	return fs
}

func DropFluxSink(s cesium.Subscriber, c cesium.Canceller) *FluxSink {
	requestedMux := sync.Mutex{}
	requested := int64(0)
	unbounded := false

	fs := &FluxSink{
		next: func(t cesium.T) {
			if c.IsCancelled() {
				return
			}

			requestedMux.Lock()
			if unbounded {
				requestedMux.Unlock()
				s.OnNext(t)
				return
			}

			if requested > 0 {
				requested = requested - 1
				requestedMux.Unlock()
				s.OnNext(t)
				return
			}
			requestedMux.Unlock()
		},
		complete: func() {
			if c.IsCancelled() {
				return
			}

			s.OnComplete()
		},
		error: func(err error) {
			if c.IsCancelled() {
				return
			}

			s.OnError(err)
		},
		internalRequest: func(n int64) {
			if c.IsCancelled() {
				return
			}

			requestedMux.Lock()
			if n == math.MaxInt64 {
				unbounded = true
				requestedMux.Unlock()
				return
			}

			requested = requested + n
			requestedMux.Unlock()
		},
		canceller: c,
	}

	c.OnCancel(fs.Cancel)

	return fs
}

func ErrorFluxSink(s cesium.Subscriber, c cesium.Canceller) *FluxSink {
	requestedMux := sync.Mutex{}
	requested := int64(0)

	closed := false
	closedMux := sync.Mutex{}
	unbounded := false

	fs := &FluxSink{
		next: func(t cesium.T) {
			if c.IsCancelled() {
				return
			}

			closedMux.Lock()
			if closed {
				closedMux.Unlock()
				return
			}
			closedMux.Unlock()

			requestedMux.Lock()
			if unbounded {
				requestedMux.Unlock()
				s.OnNext(t)
				return
			}

			if requested > 0 {
				requested = requested - 1
				requestedMux.Unlock()
				s.OnNext(t)
				return
			}
			requestedMux.Unlock()

			closedMux.Lock()
			closed = true
			closedMux.Unlock()
			s.OnError(cesium.DownstreamUnableToKeepUpError)
		},
		complete: func() {
			if c.IsCancelled() {
				return
			}

			closedMux.Lock()
			if closed {
				closedMux.Unlock()
				return
			}
			closedMux.Unlock()

			s.OnComplete()
		},
		error: func(err error) {
			if c.IsCancelled() {
				return
			}

			closedMux.Lock()
			if closed {
				closedMux.Unlock()
				return
			}
			closedMux.Unlock()

			requestedMux.Lock()
			if requested > 0 {
				requested = requested - 1
				requestedMux.Unlock()
				s.OnError(err)
				return
			}
			requestedMux.Unlock()

		},
		internalRequest: func(n int64) {
			if c.IsCancelled() {
				return
			}

			requestedMux.Lock()
			if n == math.MaxInt64 {
				unbounded = true
				requestedMux.Unlock()
				return
			}

			requested = requested + n
			requestedMux.Unlock()
		},
		canceller: c,
	}

	c.OnCancel(fs.Cancel)

	return fs
}

func IgnoreFluxSink(s cesium.Subscriber, c cesium.Canceller) *FluxSink {
	fs := &FluxSink{
		next: func(t cesium.T) {
			if c.IsCancelled() {
				return
			}

			s.OnNext(t)
		},
		complete: func() {
			if c.IsCancelled() {
				return
			}

			s.OnComplete()
		},
		error: func(err error) {
			if c.IsCancelled() {
				return
			}

			s.OnError(err)
		},
		internalRequest: func(n int64) {
		},
		canceller: c,
	}

	c.OnCancel(fs.Cancel)

	return fs
}

type SynchronousSink struct {
	emitted    bool
	emittedMux sync.Mutex
	emission   cesium.MaterializedEmission
}

func (s *SynchronousSink) Next(t cesium.T) {
	s.emittedMux.Lock()
	if s.emitted {
		s.emittedMux.Unlock()
		return
	}
	s.emitted = true
	s.emittedMux.Unlock()

	s.emission = cesium.MaterializedEmission{EventType: "next", Value: t}
}

func (s *SynchronousSink) Complete() {
	s.emittedMux.Lock()
	if s.emitted {
		s.emittedMux.Unlock()
		return
	}
	s.emitted = true
	s.emittedMux.Unlock()

	s.emission = cesium.MaterializedEmission{EventType: "complete"}
}

func (s *SynchronousSink) Error(err error) {
	s.emittedMux.Lock()
	if s.emitted {
		s.emittedMux.Unlock()
		return
	}
	s.emitted = true
	s.emittedMux.Unlock()

	s.emission = cesium.MaterializedEmission{EventType: "error", Err: err}
}

func (s *SynchronousSink) GetEmission() cesium.MaterializedEmission {
	return s.emission
}

type MonoSink struct {
	onDisposeMux sync.Mutex
	onDispose    func()

	cancelerMux sync.Mutex
	canceller   cesium.Canceller
	onCancel    func()

	completeWith    func(cesium.T)
	complete        func()
	error           func(error)
	internalRequest func(int64)

	closedMux sync.Mutex
	closed    bool
}

func (f *MonoSink) CompleteWith(t cesium.T) {
	f.closedMux.Lock()
	if f.closed {
		f.closedMux.Unlock()
		return
	}
	f.closed = true
	f.closedMux.Unlock()

	f.completeWith(t)

	f.onDisposeMux.Lock()
	if f.onDispose != nil {
		f.onDispose()
	}
	f.onDisposeMux.Unlock()
}

func (f *MonoSink) Complete() {
	f.closedMux.Lock()
	if f.closed {
		f.closedMux.Unlock()
		return
	}
	f.closed = true
	f.closedMux.Unlock()

	f.complete()

	f.onDisposeMux.Lock()
	if f.onDispose != nil {
		f.onDispose()
	}
	f.onDisposeMux.Unlock()
}

func (f *MonoSink) Error(err error) {
	f.closedMux.Lock()
	if f.closed {
		f.closedMux.Unlock()
		return
	}
	f.closed = true
	f.closedMux.Unlock()

	f.error(err)

	f.onDisposeMux.Lock()
	if f.onDispose != nil {
		f.onDispose()
	}
	f.onDisposeMux.Unlock()
}

func (f *MonoSink) OnDispose(fn func()) {
	f.onDisposeMux.Lock()
	f.onDispose = fn

	f.closedMux.Lock()
	if f.closed {
		f.onDispose()
	}
	f.closedMux.Unlock()
	f.onDisposeMux.Unlock()
}

func (f *MonoSink) OnCancel(fn func()) {
	onCancel := func() {
		f.closedMux.Lock()
		f.closed = true
		f.closedMux.Unlock()
		f.Dispose()
		fn()
	}

	f.cancelerMux.Lock()
	if f.canceller.IsCancelled() {
		onCancel()
	}
	f.onCancel = onCancel
	f.cancelerMux.Unlock()
}

func (f *MonoSink) Cancel() {
	f.cancelerMux.Lock()
	if f.onCancel != nil {
		f.onCancel()
	}
	f.cancelerMux.Unlock()

}

func (f *MonoSink) Dispose() {
	f.onDisposeMux.Lock()
	if f.onDispose != nil {
		f.onDispose()
	}
	f.onDisposeMux.Unlock()
}

func (f *MonoSink) Request(n int64) {
	f.internalRequest(n)
}

func BufferMonoSink(s cesium.Subscriber, c cesium.Canceller) *MonoSink {
	mux := sync.Mutex{}
	buffer := []cesium.MaterializedEmission{}
	requested := false

	ms := &MonoSink{
		completeWith: func(t cesium.T) {
			if c.IsCancelled() {
				return
			}

			mux.Lock()
			if requested {
				mux.Unlock()

				if c.IsCancelled() {
					return
				}

				s.OnNext(t)
				s.OnComplete()
				return
			}

			buffer = append(buffer, cesium.MaterializedEmission{EventType: "next", Value: t}, cesium.MaterializedEmission{EventType: "complete"})
			mux.Unlock()
		},
		complete: func() {
			if c.IsCancelled() {
				return
			}

			mux.Lock()
			if requested {
				mux.Unlock()

				if c.IsCancelled() {
					return
				}

				s.OnComplete()
				return
			}

			buffer = append(buffer, cesium.MaterializedEmission{EventType: "complete"})
			mux.Unlock()
		},
		error: func(err error) {
			if c.IsCancelled() {
				return
			}

			mux.Lock()
			if requested {
				mux.Unlock()

				if c.IsCancelled() {
					return
				}

				s.OnError(err)
				return
			}

			buffer = append(buffer, cesium.MaterializedEmission{EventType: "error", Err: err})
			mux.Unlock()
		},
		internalRequest: func(n int64) {
			mux.Lock()
			if c.IsCancelled() {
				mux.Unlock()
				return
			}

			for len(buffer) > 0 {
				emission := buffer[0]
				buffer = buffer[1:]

				switch emission.EventType {
				case "next":
					if c.IsCancelled() {
						mux.Unlock()
						return
					}
					s.OnNext(emission.Value)
				case "complete":
					if c.IsCancelled() {
						mux.Unlock()
						return
					}
					s.OnComplete()
				case "error":
					if c.IsCancelled() {
						mux.Unlock()
						return
					}
					s.OnError(emission.Err)
				}
			}
			requested = true
			mux.Unlock()
		},
		canceller: c,
	}

	c.OnCancel(ms.Cancel)

	return ms
}
