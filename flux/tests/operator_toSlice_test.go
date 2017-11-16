package tests

import (
	"testing"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
)

func TestToSlice(t *testing.T) {
	s, err := flux.
		Just(1, 2, 3).
		ToSlice()

	if err != nil {
		t.Errorf("Unexpected error returned. %v", err)
		return
	}

	x := []cesium.T{1, 2, 3}
	if len(s) != 3 {
		t.Errorf("Invalid output. Expected: %v, Got: %v", x, s)
		return
	}

	for i := range s {
		if s[i] != x[i] {
			t.Errorf("Invalid output. Expected: %v, Got: %v", x, s)
			return
		}
	}
}
