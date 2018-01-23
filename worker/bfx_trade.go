package worker

import (
	"bfxtrader/job"
	"bfxtrader/model"
	"fmt"
)

func opposite(side string) string {
	if side == "BUY" {
		return "SELL"
	}
	if side == "SELL" {
		return "BUY"
	}
	fmt.Println("Error int Opposite")
	return "ERROR"
}

//BFXTrade Trade BFX
func BFXTrade(market *model.BFXMarket, finish chan bool) {
	size := float32(0.01)
	currentSize := float32(0)
	currentSide := market.GetCurrentSignal()
	fmt.Println("MarketOrder   side:", currentSide, ", size:", size)
	job.MarketOrder(currentSide, size)
	for <-finish {
		job.CancelAllOrder()
		currentSize, currentSide = job.GetCurrentBFX()
		if currentSide == "" {
			fmt.Println("No Postion")
			currentSide = market.GetCurrentSignal()
			fmt.Println("MarketOrder   side:", currentSide, ", size:", size)
			job.MarketOrder(currentSide, size)
		}
		nextPrice, nextSide := market.CalcNextCross()
		if nextSide == currentSide {

			fmt.Println("MarketOrder   side:", opposite(currentSide), ", size:", currentSize)
			job.MarketOrder(opposite(currentSide), currentSize)
			fmt.Println("MarketOrder   side:", opposite(currentSide), ", size:", size)
			job.MarketOrder(opposite(currentSide), size)
			currentSide = opposite(currentSide)
		} else {
			if currentSize > size+0.005 {
				fmt.Println("MarketOrder   side:", opposite(currentSide), ", size:", currentSize-size)
				job.MarketOrder(opposite(currentSide), currentSize-size)
			}
			if currentSize < size-0.005 {
				fmt.Println("MarketOrder   side:", currentSide, ", size:", size-currentSize)
				job.MarketOrder(currentSide, size-currentSize)
			}
		}
		fmt.Println("IFDOrder  side:", opposite(currentSide), ", trigger:", nextPrice, ", size:", size)
		job.IFDStopMarketOrder(opposite(currentSide), nextPrice, size)
		collateral := job.GetCollateralValues()
		fmt.Println("Current Collateral: ", collateral.Collateral)
	}
}
