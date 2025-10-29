package kline

import (
	"fmt"

	"github.com/adshao/go-binance/v2/futures"
)

func ParseContinuousKlineEvent(event *futures.WsContinuousKlineEvent) {
}

func ErrHandler(err error) {
	fmt.Println(err)
}

func WsConnect() {
	var subscribeArgsList []*futures.WsContinuousKlineSubscribeArgs
	subscribeArgsList = append(subscribeArgsList, &futures.WsContinuousKlineSubscribeArgs{
		Pair:     "BTCUSDT",
		Interval: "1m",
	})
	doneC, _, err := futures.WsCombinedContinuousKlineServe(subscribeArgsList, ParseContinuousKlineEvent, ErrHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	<-doneC
}
