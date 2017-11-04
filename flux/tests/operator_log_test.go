package tests

import (
	"testing"

	"bytes"
	"log"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/verifier"
	"github.com/pkg/errors"
)

func TestLog(t *testing.T) {
	buffer := new(bytes.Buffer)
	expectBuffer := func(s string) {
		val := buffer.String()
		if val != s {
			t.Errorf("Wrong log output. Expected: %v, Got: %v", s, val)
		}
	}

	publisher := flux.
		FromSlice([]cesium.T{1, 2, 3}).
		Log(log.New(buffer, "", 0))

	verifier.
		Create(publisher).
		Then(func() {
			expectBuffer("Subscribed\n")
		}).
		ExpectNext(1).
		Then(func() {
			expectBuffer("Subscribed\nRequest: 1\nNext: 1\n")
		}).
		ThenRequest(2).
		ExpectNext(2, 3).
		ExpectComplete().
		Then(func() {
			expectBuffer("Subscribed\nRequest: 1\nNext: 1\nRequest: 2\nNext: 2\nNext: 3\nComplete\n")
		}).
		Verify(t)
}

func TestLogWithError(t *testing.T) {
	buffer := new(bytes.Buffer)
	expectBuffer := func(s string) {
		val := buffer.String()
		if val != s {
			t.Errorf("Wrong log output. Expected: %v, Got: %v", s, val)
		}
	}

	err := errors.New("err")

	publisher := flux.
		Error(err).
		Log(log.New(buffer, "", 0))

	verifier.
		Create(publisher).
		ExpectError(err).
		Then(func() {
			expectBuffer("Subscribed\nError: err\n")
		}).
		Verify(t)
}

func TestLogCancel(t *testing.T) {
	buffer := new(bytes.Buffer)
	expectBuffer := func(s string) {
		val := buffer.String()
		if val != s {
			t.Errorf("Wrong log output. Expected: %v, Got: %v", s, val)
		}
	}

	publisher := flux.
		Just(1).
		Log(log.New(buffer, "", 0))

	verifier.
		Create(publisher).
		ThenCancel().
		Then(func() {
			expectBuffer("Subscribed\nCancel\n")
		}).
		Verify(t)
}
