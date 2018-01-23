package worker

import (
	"bfxtrader/job"
	"bfxtrader/model"
	"sync"
	"time"
)

//BFXCalc Get BFX price and calc market
func BFXCalc(market *model.BFXMarket, finish chan bool) {
	sleepTime := time.Duration(market.SecInterval()) * time.Second
	nextTime := market.GetLastValues().GetTime().Add(sleepTime)
	time.Sleep(time.Until(nextTime))
	ch := make(chan *model.BFXPrice)
	m := new(sync.Mutex)
	for {
		go job.GetBFXPrice(ch)
		go job.CalcBFXData(ch, finish, m, market)
		time.Sleep(time.Minute)
	}
}
