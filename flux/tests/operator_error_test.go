package tests

import (
	"testing"

	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
	"github.com/pkg/errors"
)

func TestError(t *testing.T) {
	err := errors.New("err")

	verifier.
		Create(flux.Error(err)).
		ExpectError(err).
		Verify(t)
}
