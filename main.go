package main

import (
	"bfxtrader/model"
	"bfxtrader/service"
	"bfxtrader/worker"
)

func main() {
	lineClient := service.NewLineClient()
	lineClient.PushTextMessage("Program Started!")
	bfxMarket := model.NewBFXMarket(100, 21, 10, 60)
	ch := make(chan bool)
	go worker.BFXCalc(bfxMarket, ch)
	worker.BFXTrade(bfxMarket, ch)
}
