package main

import (
	"fmt"

	"time"

	"github.com/DusanKasan/cesium"
	"github.com/DusanKasan/cesium/mono"
)

type dummySubscriber struct {
}

func (dummySubscriber) OnNext(t cesium.T) {
	fmt.Println("OnNext", t)
}

func (dummySubscriber) OnError(err error) {
	fmt.Println("OnError", err)
}

func (dummySubscriber) OnComplete() {
	fmt.Println("OnComplete")
}

func (dummySubscriber) OnSubscribe(s cesium.Subscription) {
	fmt.Println("OnSubscribe", s)
}

func main() {
	mono.
		Just(1).
		Filter(func(t cesium.T) bool {
			return t == 1
		}).
		Map(func(t cesium.T) cesium.T {
			return t.(int) + 1
		}).Subscribe(dummySubscriber{}).Request(1)

	time.Sleep(time.Second)
}

//func handle(request cesium.Mono /*<Request>*/) cesium.Mono /*Response*/ {
//
//}
