package tests

import (
	"testing"

	"time"

	"errors"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
)

func TestBlockLastTimeout(t *testing.T) {
	val, ok, err := flux.
		Just(7).
		BlockLastTimeout(time.Second)

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

func TestBlockLastTimeoutWithError(t *testing.T) {
	e := errors.New("err")

	val, ok, err := flux.
		Error(e).
		BlockLastTimeout(time.Second)

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

func TestBlockLastTimeoutWithEmpty(t *testing.T) {
	val, ok, err := flux.
		Empty().
		BlockLastTimeout(time.Second)

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

func TestBlockLastTimeoutWithTimeout(t *testing.T) {
	val, ok, err := flux.
		Never().
		BlockLastTimeout(time.Millisecond)

	if val != nil {
		t.Errorf("Wrong value received. Expected: %v, Got: %v", nil, val)
	}

	if ok {
		t.Errorf("Wrong ok received. Expected: %v, Got: %v", false, ok)
	}

	if err == nil {
		t.Errorf("Wrong err received. Expected: %v, Got: %v", cesium.TimeoutError, err)
	}
}
