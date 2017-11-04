package tests

import (
	"testing"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/mono"
	"github.com/DusanKasan/cesium/verifier"
)

func TestDoOnNext(t *testing.T) {
	var buf []cesium.T

	publisher := mono.
		Just(1).
		DoOnNext(func(i cesium.T) {
			buf = append(buf, i)
		})

	verifier.
		Create(publisher).
		ExpectNext(1).
		Then(func() {
			if len(buf) != 1 || buf[0] != 1 {
				t.Errorf("Wrong data in buffer. Expected: 1, Got: %v", buf)
			}
		}).
		ExpectComplete().
		Verify(t)

}
