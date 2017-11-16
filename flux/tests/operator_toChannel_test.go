package tests

import (
	"testing"

	"errors"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
)

func TestToChannel(t *testing.T) {
	items, errs := flux.
		Just(1, 2, 3).
		ToChannel()

	var is []cesium.T
	var es []error

loop:
	for {
		select {
		case i, ok := <-items:
			if !ok {
				break loop
			}
			is = append(is, i)
		case err, ok := <-errs:
			if !ok {
				break loop
			}
			es = append(es, err)
		}
	}

	if len(es) > 0 {
		t.Errorf("error received: %v", es[0])
	}

	x := []cesium.T{1, 2, 3}
	if len(is) != 3 {
		t.Errorf("Invalid output. Expected: %v, Got: %v", x, is)
		return
	}

	for i := range is {
		if is[i] != x[i] {
			t.Errorf("Invalid output. Expected: %v, Got: %v", x, is)
			return
		}
	}
}

func TestToChanneWithError(t *testing.T) {
	originalErr := errors.New("err")
	items, errs := flux.
		Create(func(sink cesium.FluxSink) {
			sink.Next(1)
			sink.Error(originalErr)
		}, flux.OverflowStrategyBuffer).
		ToChannel()

	var is []cesium.T
	var es []error

loop:
	for {
		select {
		case i, ok := <-items:
			if !ok {
				break loop
			}
			is = append(is, i)
		case err, ok := <-errs:
			if !ok {
				break loop
			}
			es = append(es, err)
		}
	}

	if len(es) != 1 {
		t.Errorf("wrong number of errors received: %v. Details %v", len(es), es)
		return
	}

	if es[0] != originalErr {
		t.Errorf("Wrong error received. Expected %v, Got: %v", originalErr, es[0])
		return
	}

	x := []cesium.T{1}
	if len(is) != 1 {
		t.Errorf("Invalid output. Expected: %v, Got: %v", x, is)
		return
	}

	for i := range is {
		if is[i] != x[i] {
			t.Errorf("Invalid output. Expected: %v, Got: %v", x, is)
			return
		}
	}
}
