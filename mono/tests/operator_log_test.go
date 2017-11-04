package tests

import (
	"testing"

	"bytes"
	"log"

	"github.com/DusanKasan/cesium/mono"
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

	publisher := mono.
		Just(1).
		Log(log.New(buffer, "", 0))

	verifier.
		Create(publisher).
		Then(func() {
			expectBuffer("Subscribed\n")
		}).
		ExpectNext(1).
		ExpectComplete().
		Then(func() {
			expectBuffer("Subscribed\nRequest: 1\nNext: 1\nComplete\n")
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

	publisher := mono.
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

	publisher := mono.
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
