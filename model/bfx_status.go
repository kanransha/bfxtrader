package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"../service"
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
	url := "/v1/me/getexecutions?product_code=FX_BTC_JPY&count=1"
	if lastID != 0 {
		url = "/v1/me/getexecutions?product_code=FX_BTC_JPY&after=" + fmt.Sprint(lastID)
	}
	client := service.NewBitClient()
	request, err := client.NewRequest(url, "GET", "")
	if err != nil {
		fmt.Println("getExecutions Request Error")
		return getExecutions(lastID)
	}
	fmt.Println("Req", request)
	res, err := client.Do(request)
	if err != nil {
		fmt.Println("getExecutions Response Error")
		return getExecutions(lastID)
	}
	jsonBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("getExecutions Read Body Error")
		res.Body.Close()
		return getExecutions(lastID)
	}
	fmt.Println("Res", string(jsonBytes))
	if res.StatusCode != 200 {
		fmt.Printf("getExecutions StatusCode = %d\n", res.StatusCode)
		fmt.Println(string(jsonBytes))
		res.Body.Close()
		return getExecutions(lastID)
	}
	defer res.Body.Close()
	jsonData := new(executions)
	if err := json.Unmarshal(jsonBytes, jsonData); err != nil {
		fmt.Println("getExecutions JSON Unmarchal error")
		fmt.Println(string(jsonBytes))
		return getExecutions(lastID)
	}
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
