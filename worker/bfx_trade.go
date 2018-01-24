package worker

import (
	"fmt"
	"fxtrader/job"
	"fxtrader/model"
	"fxtrader/mytype"
)

func updateBFXPosition(side *mytype.Side, size float32) {
	currentSize, sideStr := job.GetCurrentBFX()
	currentSide := mytype.NewSide(sideStr)
	fmt.Println("Position Now:  ", currentSide.String(), currentSize, "BFX")
	fmt.Println("Position Next: ", side.String(), size, "BFX")
	if currentSide.IsZero() {
		if side.IsZero() {
			return
		}
		job.BFXMarketOrder(side.String(), size)
		return
	}
	if side.IsSame(currentSide) {
		if size < currentSize-0.001 {
			job.BFXMarketOrder(side.Opposite().String(), currentSize-size)
			return
		}
		if size > currentSize+0.001 {
			job.BFXMarketOrder(side.String(), size-currentSize)
			return
		}
		return
	}
	if side.IsOpposite(currentSide) {
		job.BFXMarketOrder(side.String(), size)
		job.BFXMarketOrder(side.String(), currentSize)
		return
	}
	job.BFXMarketOrder(currentSide.Opposite().String(), currentSize)
}

//BFXTrade Trade BFX
func BFXTrade(market *model.BFXMarket, finish chan bool) {
	size := float32(0.01)
	side := mytype.NewSide(market.GetCurrentSignal())
	updateBFXPosition(side, size)
	for <-finish {
		job.CancelAllBFXOrder()
		side := mytype.NewSide(market.GetCurrentSignal())
		updateBFXPosition(side, size)
		nextPrice, nextSideStr := market.CalcNextCross()
		nextSide := mytype.NewSide(nextSideStr)
		job.BFXIFDStopMarketOrder(nextSide.String(), nextPrice, size)
		collateral := job.GetBFXCollateralValues()
		fmt.Println("Current Collateral: ", collateral.Collateral)
	}
}
