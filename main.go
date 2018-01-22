package main

import (
	"./model"
	"./worker"
)

func main() {
	bfxMarket := model.NewBFXMarket(100, 21, 10, 60)
	ch := make(chan bool)
	go worker.BFXCalc(bfxMarket, ch)
	worker.BFXTrade(bfxMarket, ch)
}
