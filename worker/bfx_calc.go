package worker

import (
	"time"

	"../job"
	"../model"
)

//BFXCalc Get BFX price and calc market
func BFXCalc(market *model.BFXMarket) {
	sleepTime := time.Duration(market.SecInterval()) * time.Second
	nextTime := market.GetLastValues().GetTime().Add(sleepTime)
	time.Sleep(time.Until(nextTime))
	ch := make(chan *model.BFXPrice)
	for {
		go job.GetBFXPrice(ch)
		go job.CalcBFXData(ch, market)
		time.Sleep(time.Minute)
	}
}
