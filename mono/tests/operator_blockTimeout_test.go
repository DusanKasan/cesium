package tests

import (
	"testing"

	"time"

	"errors"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/mono"
)

func TestBlockTimeout(t *testing.T) {
	val, ok, err := mono.
		Just(5).
		BlockTimeout(time.Second)

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

func TestBlockTimeoutWithError(t *testing.T) {
	e := errors.New("err")

	val, ok, err := mono.
		Error(e).
		BlockTimeout(time.Second)

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

func TestBlockTimeoutWithEmpty(t *testing.T) {
	val, ok, err := mono.
		Empty().
		BlockTimeout(time.Second)

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

func TestBlockTimeoutWithTimeout(t *testing.T) {
	val, ok, err := mono.
		Never().
		BlockTimeout(time.Millisecond)

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
