package job

import (
	"fmt"
	"fxtrader/service"
)

type bfxParentOrderParam struct {
	ProductCode   string  `json:"product_code"`
	ConditionType string  `json:"condition_type"`
	Side          string  `json:"side"`
	Price         float32 `json:"price"`
	TriggerPrice  float32 `json:"trigger_price"`
	Size          float32 `json:"size"`
}

type bfxParentOrderBody struct {
	OrderMethod    string                `json:"order_method"`
	MinuteToExpire int                   `json:"minute_to_expire"`
	TimeInForce    string                `json:"time_in_force"`
	Parameters     []bfxParentOrderParam `json:"parameters"`
}

//BFXParentOrderAcceptanceID Parent order acceptance ID
type BFXParentOrderAcceptanceID string

type bfxParentOrderResponseBody struct {
	ID BFXParentOrderAcceptanceID `json:"parent_order_acceptance_id"`
}

func newBFXParentOrderParam(productCode string, conditionType string, side string, price float32, triggerPrice float32, size float32) *bfxParentOrderParam {
	return &bfxParentOrderParam{productCode, conditionType, side, price, triggerPrice, size}
}

func bfxParentOrder(orderMethod string, minuteToExpire int, timeInForce string, parameters []bfxParentOrderParam) BFXParentOrderAcceptanceID {
	postBody := bfxParentOrderBody{orderMethod, minuteToExpire, timeInForce, parameters}
	client := service.NewBitClient()
	pathDir := "/v1/me/sendparentorder"
	response := new(bfxParentOrderResponseBody)
	client.Post(pathDir, postBody, response)
	return response.ID
}

//BFXIFDStopMarketOrder Do parent order of IFD, STOP order first and then MARKET order of same size
func BFXIFDStopMarketOrder(side string, triggerPrice float32, size float32) BFXParentOrderAcceptanceID {
	fmt.Printf("IFDOrder  side: %s, trigger: %.0f , size: %f\n", side, triggerPrice, size)
	productType := "FX_BTC_JPY"
	firstOrder := newBFXParentOrderParam(productType, "STOP", side, 0, triggerPrice, size)
	secondOrder := newBFXParentOrderParam(productType, "MARKET", side, 0, 0, size)
	parameters := []bfxParentOrderParam{*firstOrder, *secondOrder}
	return bfxParentOrder("IFD", 2, "GTC", parameters)
}

//BFXStopOrder Do parent order of STOP
func BFXStopOrder(side string, triggerPrice float32, size float32) BFXParentOrderAcceptanceID {
	order := newBFXParentOrderParam("FX_BTC_JPY", "STOP", side, 0, triggerPrice, size)
	return bfxParentOrder("SIMPLE", 2, "GTC", []bfxParentOrderParam{*order})
}
