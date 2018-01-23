package model

import (
	"bfxtrader/service"
	"fmt"
)

//BFXStatus Store BFX status
type BFXStatus struct {
	lastID int
	size   float32
	side   string
}

type execution struct {
	ID         int     `json:"id"`
	OrderID    string  `json:"child_order_id"`
	Side       string  `json:"side"`
	Price      float32 `json:"price"`
	Size       float32 `json:"size"`
	Commission float32 `json:"commission"`
	Date       string  `json:"exec_date"`
	AcceptID   string  `json:"child_order_acceptance_id"`
}

type executions []execution

func getExecutions(lastID int) *executions {
	pathDir := "/v1/me/getexecutions"
	queryStr := "product_code=FX_BTC_JPY&count=1"
	if lastID != 0 {
		queryStr = "product_code=FX_BTC_JPY&after=" + fmt.Sprint(lastID)
	}
	client := service.NewBitClient()
	jsonData := new(executions)
	client.Get(pathDir, queryStr, jsonData)
	return jsonData
}

//NewBFXStatus Create new BFX status
func NewBFXStatus() *BFXStatus {
	exes := getExecutions(0)
	status := new(BFXStatus)
	status.lastID = (*exes)[0].ID
	status.side = "BUY"
	status.size = 0
	return status
}
