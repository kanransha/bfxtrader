package job

import (
	"../model"
)

//CalcBFXData Calculate BFX Data
func CalcBFXData(c chan *model.BFXPrice, market *model.BFXMarket) {
	currentPrice := <-c
	market.SetBFXPrice(currentPrice)
}
