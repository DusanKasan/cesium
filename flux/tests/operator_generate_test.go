package tests

import (
	"testing"

	"math"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
	"github.com/pkg/errors"
)

func TestGenerate(t *testing.T) {
	i := 1

	publisher := flux.Generate(func(s cesium.SynchronousSink) {
		if i == 3 {
			// this also tests that only the first invocation on
			// cesium.SynchronousSink is accepted
			s.Complete()
		}

		s.Next(i)
		i++
	})

	verifier.
		Create(publisher).
		ExpectNext(1, 2).
		ThenRequest(1).
		ExpectComplete().
		Verify(t)
}

func TestGenerateSubscriptionBehaviour(t *testing.T) {
	i := 1

	publisher := flux.Generate(func(s cesium.SynchronousSink) {
		if i == 3 {
			// this also tests that only the first invocation on
			// cesium.SynchronousSink is accepted
			s.Complete()
		}

		s.Next(i)
		i++
	})

	verifier.
		Create(publisher).
		ExpectNext(1, 2).
		ThenCancel().
		ThenRequest(1).
		ExpectNextCount(0).
		Verify(t)
}

func TestGenerateOnError(t *testing.T) {
	i := 1
	err := errors.New("err")

	publisher := flux.Generate(func(s cesium.SynchronousSink) {
		if i == 3 {
			// this also tests that only the first invocation on
			// cesium.SynchronousSink is accepted
			s.Error(err)
		}

		s.Next(i)
		i++
	})

	verifier.
		Create(publisher).
		ExpectNext(1, 2).
		ThenRequest(1).
		ExpectError(err).
		Verify(t)
}

func TestGenerateUnbounded(t *testing.T) {
	i := 1

	publisher := flux.Generate(func(s cesium.SynchronousSink) {
		if i == 3 {
			// this also tests that only the first invocation on
			// cesium.SynchronousSink is accepted
			s.Complete()
		}

		s.Next(i)
		i++
	})

	verifier.
		Create(publisher).
		ThenRequest(math.MaxInt64).
		ExpectNextCount(2).
		ExpectComplete().
		Verify(t)
}

func TestGenerateWithNoEmissions(t *testing.T) {
	i := 1

	publisher := flux.Generate(func(s cesium.SynchronousSink) {
		if i == 3 {
			s.Complete()
		}

		i++
	})

	verifier.
		Create(publisher).
		ThenRequest(1).
		ExpectError(cesium.NoEmissionOnSynchronousSinkError).
		Verify(t)
}
