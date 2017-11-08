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
