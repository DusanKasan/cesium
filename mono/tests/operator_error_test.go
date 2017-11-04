package tests

import (
	"testing"

	"github.com/DusanKasan/cesium/mono"
	"github.com/DusanKasan/cesium/verifier"
	"github.com/pkg/errors"
)

func TestError(t *testing.T) {
	err := errors.New("err")

	verifier.
		Create(mono.Error(err)).
		ExpectError(err).
		Verify(t)
}
