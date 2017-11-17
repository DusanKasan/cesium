package tests

import (
	"testing"

	"errors"

	"github.com/DusanKasan/cesium/mono"
	"github.com/DusanKasan/cesium/verifier"
)

func TestOnErrorResume(t *testing.T) {
	err := errors.New("err")

	f := mono.
		Error(err).
		OnErrorResume(
			func(e error) bool {
				return e == err
			},
			mono.Just(1))

	verifier.
		Create(f).
		ExpectNext(1).
		ExpectComplete().
		Verify(t)
}
