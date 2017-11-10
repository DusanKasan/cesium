package json

import (
	"encoding/json"
	"net/http"

	"io"
	"strings"

	"github.com/DusanKasan/cesium"
)

func EncodeResponseFromMono(m cesium.Mono, w http.ResponseWriter) cesium.Mono {
	return m.Handle(func(t cesium.T, sink cesium.SynchronousSink) {
		v, err := json.Marshal(t)
		if err != nil {
			sink.Error(err)
			return
		}

		_, err = io.Copy(w, strings.NewReader(string(v)))
		if err != nil {
			sink.Error(err)
			return
		}

		sink.Next(t)
	})
}
