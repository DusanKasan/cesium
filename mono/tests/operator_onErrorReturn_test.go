package tests

import (
	"testing"

	"github.com/DusanKasan/cesium/mono"
	"github.com/DusanKasan/cesium/verifier"
	"github.com/pkg/errors"
)

func TestOnErrorReturn(t *testing.T) {
	publisher := mono.
		Error(errors.New("err")).
		OnErrorReturn(2)

	verifier.Create(publisher).
		ExpectNext(2).
		ExpectComplete().
		Verify(t)
}
