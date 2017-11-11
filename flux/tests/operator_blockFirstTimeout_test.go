package tests

import (
	"testing"

	"time"

	"errors"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
)

func TestBlockFirstTimeout(t *testing.T) {
	val, ok, err := flux.
		Just(5).
		BlockFirstTimeout(time.Second)

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

func TestBlockFirstTimeoutWithError(t *testing.T) {
	e := errors.New("err")

	val, ok, err := flux.
		Error(e).
		BlockFirstTimeout(time.Second)

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

func TestBlockFirstTimeoutWithEmpty(t *testing.T) {
	val, ok, err := flux.
		Empty().
		BlockFirstTimeout(time.Second)

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

func TestBlockFirstTimeoutWithTimeout(t *testing.T) {
	val, ok, err := flux.
		Never().
		BlockFirstTimeout(time.Millisecond)

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
