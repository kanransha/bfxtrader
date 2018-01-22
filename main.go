package main

import (
	"./model"
	"./worker"
)

func main() {
	bfxMarket := model.NewBFXMarket(100, 21, 10, 60)
	worker.BFXCalc(bfxMarket)
}
