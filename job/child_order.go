package job

import (
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
	postBody := childOrderBody{productCode, orderType, side, price, size, minuteToExpire, timeInForce}
	client := service.NewBitClient()
	pathDir := "/v1/me/sendchildorder"
	response := new(childOrderResponseBody)
	client.Post(pathDir, postBody, response)
	return response.ID
}

//MarketOrder Do simple market order
func MarketOrder(side string, size float32) ChildOrderAcceptanceID {
	return childOrder("FX_BTC_JPY", "MARKET", side, 0, size, 2, "GTC")
}
