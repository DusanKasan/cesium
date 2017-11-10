package internal

import "github.com/DusanKasan/cesium"

type ScalarMono struct {
	cesium.Mono
	get func() (cesium.T, bool)
}

func (s *ScalarMono) Get() (cesium.T, bool) {
	return s.get()
}

func (s *ScalarMono) Filter(f func(cesium.T) bool) cesium.Mono {
	return MonoFilterOperator(s, f)
}

func (s *ScalarMono) Map(f func(cesium.T) cesium.T) cesium.Mono {
	return MonoMapOperator(s, f)
}

func (s *ScalarMono) FlatMap(fn func(cesium.T) cesium.Mono, scheduler ...cesium.Scheduler) cesium.Mono {
	t, ok := s.Get()
	if ok {
		return MonoDefer(func() cesium.Mono { return fn(t) })
	}

	return s
}

func (s *ScalarMono) Handle(fn func(cesium.T, cesium.SynchronousSink)) cesium.Mono {
	t, ok := s.Get()

	if !ok {
		return s
	}

	sink := &SynchronousSink{}
	fn(t, sink)
	emission := sink.GetEmission()
	switch emission.EventType {
	case "next":
		return monoFromCallable(func() (cesium.T, bool) {
			return emission.Value, true
		})
	case "complete":
		return MonoEmpty()
	case "error":
		return MonoError(emission.Err)
	default:
		return MonoError(cesium.NoEmissionOnSynchronousSinkError)
	}
}

func (s *ScalarMono) FlatMapMany(fn func(cesium.T) cesium.Publisher, scheduler ...cesium.Scheduler) cesium.Flux {
	t, ok := s.Get()
	if ok {
		return FluxDefer(func() cesium.Publisher { return fn(t) })
	}

	return FluxEmpty()
}
