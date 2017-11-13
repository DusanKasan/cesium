package internal

import "github.com/DusanKasan/cesium"

type ScalarCallable interface {
	Get() (cesium.T, bool)
}

type ScalarFlux struct {
	cesium.Flux
	get func() (cesium.T, bool)
}

func (s *ScalarFlux) Get() (cesium.T, bool) {
	return s.get()
}

func (s *ScalarFlux) Count() cesium.Mono {
	return FluxCountOperator(s)
}

func (s *ScalarFlux) Filter(f func(cesium.T) bool) cesium.Flux {
	return FluxFilterOperator(s, f)
}

func (s *ScalarFlux) Map(f func(cesium.T) cesium.T) cesium.Flux {
	return FluxMapOperator(s, f)
}

func (s *ScalarFlux) Reduce(f func(cesium.T, cesium.T) cesium.T) cesium.Mono {
	return monoFromCallable(func() (cesium.T, bool) {
		return s.Get()
	})
}

func (s *ScalarFlux) Scan(f func(cesium.T, cesium.T) cesium.T) cesium.Flux {
	return s
}

func (s *ScalarFlux) All(f func(cesium.T) bool) cesium.Mono {
	return monoFromCallable(func() (cesium.T, bool) {
		t, ok := s.Get()
		if !ok {
			return nil, false
		}

		return f(t), true
	})
}

func (s *ScalarFlux) Any(f func(cesium.T) bool) cesium.Mono {
	return s.All(f)
}

func (s *ScalarFlux) HasElements() cesium.Mono {
	return monoFromCallable(func() (cesium.T, bool) {
		_, ok := s.Get()
		return ok, true
	})
}

func (s *ScalarFlux) HasElement(t cesium.T) cesium.Mono {
	return monoFromCallable(func() (cesium.T, bool) {
		tt, ok := s.Get()
		return t == tt, ok
	})
}

func (s *ScalarFlux) Handle(fn func(cesium.T, cesium.SynchronousSink)) cesium.Flux {
	t, ok := s.Get()

	if !ok {
		return s
	}

	sink := &SynchronousSink{}
	fn(t, sink)
	sig := sink.Signal()
	switch sig.Type() {
	case cesium.SignalTypeOnNext:
		return fluxFromCallable(func() (cesium.T, bool) {
			return sig.Item(), true
		})
	case cesium.SignalTypeOnComplete:
		return FluxEmpty()
	case cesium.SignalTypeOnError:
		return FluxError(sig.Error())
	default:
		return FluxError(cesium.NoEmissionOnSynchronousSinkError)
	}
}

func (s *ScalarFlux) DistinctUntilChanged() cesium.Flux {
	return s
}

func (s *ScalarFlux) FlatMap(fn func(cesium.T) cesium.Publisher, scheduler ...cesium.Scheduler) cesium.Flux {
	t, ok := s.Get()
	if ok {
		return FluxDefer(func() cesium.Publisher { return fn(t) })
	}

	return s
}
