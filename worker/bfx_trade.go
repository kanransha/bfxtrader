package worker

import (
	"bfxtrader/job"
	"bfxtrader/model"
	"fmt"
)

//Side Side BUY SELL ZERO
type Side string

//NewSide Create Side by string
func NewSide(sideStr string) *Side {
	side := Side(sideStr)
	return &side
}

//SELLSide SELL
func SELLSide() *Side {
	sideStr := "SELL"
	side := Side(sideStr)
	return &side
}

//BUYSide BUY
func BUYSide() *Side {
	sideStr := "BUY"
	side := Side(sideStr)
	return &side
}

//ZEROSide ZERO
func ZEROSide() *Side {
	sideStr := "ZERO"
	side := Side(sideStr)
	return &side
}

//IsOpposite true if opposite
func (side *Side) IsOpposite(sideInQuestion *Side) bool {
	if string(*side) == "SELL" && string(*sideInQuestion) == "BUY" {
		return true
	}
	if string(*side) == "BUY" && string(*sideInQuestion) == "SELL" {
		return true
	}
	return false
}

//IsSame true if same
func (side *Side) IsSame(sideInQuestion *Side) bool {
	if string(*side) == "SELL" && string(*sideInQuestion) == "SELL" {
		return true
	}
	if string(*side) == "BUY" && string(*sideInQuestion) == "BUY" {
		return true
	}
	if string(*side) == "ZERO" && string(*sideInQuestion) == "ZERO" {
		return true
	}
	return false
}

//Opposite Return opposite side
func (side *Side) Opposite() *Side {
	switch string(*side) {
	case "BUY":
		return SELLSide()
	case "SELL":
		return BUYSide()
	case "ZERO":
		return ZEROSide()
	}
	fmt.Println("Error in Opposite")
	return nil
}

//IsZero true if ZERO
func (side *Side) IsZero() bool {
	if string(*side) == "ZERO" {
		return true
	}
	return false
}

func (side *Side) String() string {
	return string(*side)
}

func opposite(side string) string {
	if side == "BUY" {
		return "SELL"
	}
	if side == "SELL" {
		return "BUY"
	}
	fmt.Println("Error in Opposite")
	return "ERROR"
}

func updatePosition(side *Side, size float32) {
	currentSize, sideStr := job.GetCurrentBFX()
	currentSide := NewSide(sideStr)
	fmt.Println("Position Now:  ", currentSide.String(), currentSize, "BFX")
	fmt.Println("Position Next: ", side.String(), size, "BFX")
	if currentSide.IsZero() {
		if side.IsZero() {
			return
		}
		job.MarketOrder(side.String(), size)
		return
	}
	if side.IsSame(currentSide) {
		if size < currentSize-0.001 {
			job.MarketOrder(side.Opposite().String(), currentSize-size)
			return
		}
		if size > currentSize+0.001 {
			job.MarketOrder(side.String(), size-currentSize)
			return
		}
		return
	}
	if side.IsOpposite(currentSide) {
		job.MarketOrder(side.String(), size)
		job.MarketOrder(side.String(), currentSize)
		return
	}
	job.MarketOrder(currentSide.Opposite().String(), currentSize)
}

//BFXTrade Trade BFX
func BFXTrade(market *model.BFXMarket, finish chan bool) {
	size := float32(0.01)
	side := NewSide(market.GetCurrentSignal())
	updatePosition(side, size)
	for <-finish {
		job.CancelAllOrder()
		side := NewSide(market.GetCurrentSignal())
		updatePosition(side, size)
		nextPrice, nextSideStr := market.CalcNextCross()
		nextSide := NewSide(nextSideStr)
		job.IFDStopMarketOrder(nextSide.String(), nextPrice, size)
		collateral := job.GetCollateralValues()
		fmt.Println("Current Collateral: ", collateral.Collateral)
	}
}
