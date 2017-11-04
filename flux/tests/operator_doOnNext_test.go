package tests

import (
	"testing"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
)

func TestDoOnNext(t *testing.T) {
	var buf []cesium.T

	publisher := flux.
		FromSlice([]cesium.T{1, 2, 3}).
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
		ExpectNext(2, 3).
		Then(func() {
			if len(buf) != 3 || buf[0] != 1 || buf[1] != 2 || buf[2] != 3 {
				t.Errorf("Wrong data in buffer. Expected: [1, 2, 3], Got: %v", buf)
			}
		}).
		ExpectComplete().
		Verify(t)

}
