package job

import (
	"fmt"
	"fxtrader/service"
)

type bfxChildOrderBody struct {
	ProductCode    string  `json:"product_code"`
	OrderType      string  `json:"child_order_type"`
	Side           string  `json:"side"`
	Price          float32 `json:"price"`
	Size           float32 `json:"size"`
	MinuteToExpire int     `json:"minute_to_expire"`
	TimeInForce    string  `json:"time_in_force"`
}

//BFXChildOrderAcceptanceID Parent order acceptance ID
type BFXChildOrderAcceptanceID string

type bfxChildOrderResponseBody struct {
	ID BFXChildOrderAcceptanceID `json:"child_order_acceptance_id"`
}

func bfxChildOrder(productCode string, orderType string, side string, price float32, size float32, minuteToExpire int, timeInForce string) BFXChildOrderAcceptanceID {
	postBody := bfxChildOrderBody{productCode, orderType, side, price, size, minuteToExpire, timeInForce}
	client := service.NewBitClient()
	pathDir := "/v1/me/sendchildorder"
	response := new(bfxChildOrderResponseBody)
	client.Post(pathDir, postBody, response)
	return response.ID
}

//BFXMarketOrder Do simple market order
func BFXMarketOrder(side string, size float32) BFXChildOrderAcceptanceID {
	fmt.Println("MarketOrder   side:", side, ", size:", size)
	return bfxChildOrder("FX_BTC_JPY", "MARKET", side, 0, size, 2, "FOK")
}
