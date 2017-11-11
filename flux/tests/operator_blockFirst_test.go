package tests

import (
	"testing"

	"errors"

	"github.com/DusanKasan/cesium/flux"
)

func TestBlockFirst(t *testing.T) {
	val, ok, err := flux.
		Just(5, 6, 7).
		BlockFirst()

	if val != 5 {
		t.Errorf("Wrong value received. Expected: %v, Got: %v", 5, val)
	}

	if !ok {
		t.Errorf("Wrong ok received. Expected: %v, Got: %v", true, ok)
	}

	if err != nil {
		t.Errorf("Wrong err received. Expected: %v, Got: %v", nil, err)
	}
}

func TestBlockFirstWithError(t *testing.T) {
	e := errors.New("err")

	val, ok, err := flux.
		Error(e).
		BlockFirst()

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

func TestBlockFirstWithEmpty(t *testing.T) {
	val, ok, err := flux.
		Empty().
		BlockFirst()

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
