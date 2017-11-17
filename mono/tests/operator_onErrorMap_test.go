package tests

import (
	"testing"

	"errors"

	"github.com/DusanKasan/cesium/mono"
	"github.com/DusanKasan/cesium/verifier"
)

func TestOnErrorMap(t *testing.T) {
	err := errors.New("err")
	replacementError := errors.New("repl")

	f := mono.
		Error(err).
		OnErrorMap(
			func(e error) error {
				return replacementError
			})

	verifier.
		Create(f).
		ExpectError(replacementError).
		Verify(t)
}
