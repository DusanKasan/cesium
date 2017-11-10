package tests

import (
	"testing"

	"errors"

	"github.com/DusanKasan/cesium/mono"
)

func TestBlock(t *testing.T) {
	val, ok, err := mono.
		Just(5).
		Block()

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

func TestBlockWithError(t *testing.T) {
	e := errors.New("err")

	val, ok, err := mono.
		Error(e).
		Block()

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

func TestBlockWithEmpty(t *testing.T) {
	val, ok, err := mono.
		Empty().
		Block()

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
