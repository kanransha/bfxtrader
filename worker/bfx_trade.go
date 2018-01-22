package worker

import (
	"fmt"

	"../job"
	"../model"
)

//BFXTrade Trade BFX
func BFXTrade(market *model.BFXMarket, finish chan bool) {
	size := float32(0.01)
	currentSide := market.GetCurrentSignal()
	job.MarketOrder(currentSide, size)
	for <-finish {
		job.CancelAllOrder()
		price, side := market.CalcNextCross()
		if side == currentSide {
			if side == "BUY" {
				job.MarketOrder("SELL", size)
				job.MarketOrder("SELL", size)
			}
			if side == "SELL" {
				job.MarketOrder("BUY", size)
				job.MarketOrder("BUY", size)
			}
		}
		job.IFDStopMarketOrder(side, price, size)
		collateral := job.GetCollateralValues()
		fmt.Println("Current Collateral: ", collateral.Collateral)
	}
}
