package tests

import (
	"testing"

	"github.com/DusanKasan/cesium/mono"
	"github.com/DusanKasan/cesium/verifier"
)

func TestEmpty(t *testing.T) {
	verifier.
		Create(mono.Empty()).
		ExpectComplete().
		Verify(t)
}
