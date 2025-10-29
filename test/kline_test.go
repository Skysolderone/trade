package test

import (
	"testing"

	"trade/kline"
)

func TestKline(t *testing.T) {
	kline.GetKline("BTCUSDT", "1d")
}
