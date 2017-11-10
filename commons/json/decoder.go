package json

import (
	"encoding/json"
	"io"
	"reflect"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/flux"
	"github.com/DusanKasan/cesium/mono"
)

// DecodeArray sequentially decodes elements of a json array read from a
// io.Reader into values of type i and returns a cesium.Flux of these values.
func DecodeArray(r io.Reader, i interface{}) cesium.Flux {
	inputType := reflect.TypeOf(i)
	dec := json.NewDecoder(r)
	openingBraceRead := false

	return flux.Generate(func(sink cesium.SynchronousSink) {
		if !openingBraceRead {
			openingBraceRead = true
			_, err := dec.Token()
			if err != nil {
				sink.Error(err)
				return
			}
		}

		if dec.More() {
			m := reflect.New(inputType).Interface()
			if err := dec.Decode(&m); err == io.EOF {
				sink.Complete()
			} else if err != nil {
				sink.Error(err)
			}

			sink.Next(m)
		} else {
			_, err := dec.Token()
			if err != nil {
				sink.Error(err)
			}
			sink.Complete()
		}
	})
}

// DecodeScalar decodes a value of type i and returns a cesium.Mono that emits it.
func DecodeScalar(r io.Reader, i interface{}) cesium.Mono {
	v := reflect.New(reflect.TypeOf(i)).Elem()
	x := v.Interface()
	err := json.NewDecoder(r).Decode(&x)
	if err != nil {
		return mono.Error(err)
	}

	return mono.Just(v.Interface())
}
