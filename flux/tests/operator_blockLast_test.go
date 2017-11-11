package tests

import (
	"testing"

	"errors"

	"github.com/DusanKasan/cesium/flux"
)

func TestBlockLast(t *testing.T) {
	val, ok, err := flux.
		Just(5, 6, 7).
		BlockLast()

	if val != 7 {
		t.Errorf("Wrong value received. Expected: %v, Got: %v", 7, val)
	}

	if !ok {
		t.Errorf("Wrong ok received. Expected: %v, Got: %v", true, ok)
	}

	if err != nil {
		t.Errorf("Wrong err received. Expected: %v, Got: %v", nil, err)
	}
}

func TestBlockLastWithError(t *testing.T) {
	e := errors.New("err")

	val, ok, err := flux.
		Error(e).
		BlockLast()

	if val != nil {
		t.Errorf("Wrong value received. Expected: %v, Got: %v", nil, val)
	}

	if ok {
		t.Errorf("Wrong ok received. Expected: %v, Got: %v", false, ok)
	}

	if err == nil {
		t.Errorf("Wrong err received. Expected: %v, Got: %v", e, err)
	}
}

func TestBlockLastWithEmpty(t *testing.T) {
	val, ok, err := flux.
		Empty().
		BlockLast()

	if val != nil {
		t.Errorf("Wrong value received. Expected: %v, Got: %v", nil, val)
	}

	if ok {
		t.Errorf("Wrong ok received. Expected: %v, Got: %v", false, ok)
	}

	if err != nil {
		t.Errorf("Wrong err received. Expected: %v, Got: %v", nil, err)
	}
}
