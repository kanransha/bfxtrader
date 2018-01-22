package model

import (
	"fmt"
	"time"
)

//BFXPrice BFX price with the time
type BFXPrice struct {
	endTime time.Time
	value   float32
}

//BFXPrices Type to sort BFXPrice
type BFXPrices []BFXPrice

//BFXValues BFX price and other values
type BFXValues struct {
	price   BFXPrice
	slowEMA float32
	fastEMA float32
}

func calcNextEMA(newPrice float32, lastEMA float32, period int) float32 {
	return ((float32(period)-1.0)*lastEMA + 2.0*newPrice) / (float32(period) + 1.0)
}

//NewBFXPrice Create New BFXPrice
func NewBFXPrice(endTime time.Time, value float32) *BFXPrice {
	return &BFXPrice{endTime, value}
}

func (prices BFXPrices) Len() int {
	return len(prices)
}

func (prices BFXPrices) Swap(i, j int) {
	prices[i], prices[j] = prices[j], prices[i]
}

func (prices BFXPrices) Less(i, j int) bool {
	return prices[i].endTime.Unix() < prices[j].endTime.Unix()
}

//Print Print BFXValues with JST time
func (values *BFXValues) Print() {
	location := time.FixedZone("Asia/Tokyo", 9*60*60)
	printTime := values.price.endTime.In(location)
	fmt.Printf("[%s JST]  ", printTime.Format("2006/01/02 15:04:05"))
	fmt.Printf("Price: %7.0f   Slow EMA: %7.0f   Fast EMA: %7.0f\n", values.price.value, values.slowEMA, values.fastEMA)
}

//GetTime Get time of the bfx values
func (values *BFXValues) GetTime() time.Time {
	return values.price.endTime
}
