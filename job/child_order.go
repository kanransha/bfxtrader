package job

import (
	"bytes"
	"encoding/json"
	"fmt"

	"../service"
)

type childOrderBody struct {
	ProductCode    string  `json:"product_code"`
	OrderType      string  `json:"child_order_type"`
	Side           string  `json:"side"`
	Price          float32 `json:"price"`
	Size           float32 `json:"size"`
	MinuteToExpire int     `json:"minute_to_expire"`
	TimeInForce    string  `json:"time_in_force"`
}

//ChildOrderAcceptanceID Parent order acceptance ID
type ChildOrderAcceptanceID string

type childOrderResponseBody struct {
	ID ChildOrderAcceptanceID `json:"child_order_acceptance_id"`
}

func childOrder(productCode string, orderType string, side string, price float32, size float32, minuteToExpire int, timeInForce string) ChildOrderAcceptanceID {
	jsonBytes, err := json.Marshal(childOrderBody{productCode, orderType, side, price, size, minuteToExpire, timeInForce})
	client := service.NewBitClient()
	request, err := client.NewRequest("/v1/me/sendchildorder", "POST", string(jsonBytes))
	if err != nil {
		fmt.Println("ChildOrder Request Error")
		return childOrder(productCode, orderType, side, price, size, minuteToExpire, timeInForce)
	}
	res, err := client.Do(request)
	if err != nil {
		fmt.Println("ChildOrder Response Error")
		return childOrder(productCode, orderType, side, price, size, minuteToExpire, timeInForce)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	if res.StatusCode != 200 {
		res.Body.Close()
		fmt.Printf("ChildOrder StatusCode = %d\n", res.StatusCode)
		fmt.Println(buf.String())
		return childOrder(productCode, orderType, side, price, size, minuteToExpire, timeInForce)
	}
	defer res.Body.Close()
	jsonBytesResp := buf.Bytes()
	jsonData := new(childOrderResponseBody)
	if err := json.Unmarshal(jsonBytesResp, jsonData); err != nil {
		fmt.Println("ChildOrder Response JSON Error")
		return childOrder(productCode, orderType, side, price, size, minuteToExpire, timeInForce)
	}
	defer fmt.Println("Order Completed: ", jsonData.ID)
	return jsonData.ID
}

//MarketOrder Do simple market order
func MarketOrder(side string, size float32) ChildOrderAcceptanceID {
	return childOrder("FX_BTC_JPY", "MARKET", side, 0, size, 2, "GTC")
}
