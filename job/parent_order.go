package job

import (
	"../service"
)

type parentOrderParam struct {
	ProductCode   string  `json:"product_code"`
	ConditionType string  `json:"condition_type"`
	Side          string  `json:"side"`
	Price         float32 `json:"price"`
	TriggerPrice  float32 `json:"trigger_price"`
	Size          float32 `json:"size"`
}

type parentOrderBody struct {
	OrderMethod    string             `json:"order_method"`
	MinuteToExpire int                `json:"minute_to_expire"`
	TimeInForce    string             `json:"time_in_force"`
	Parameters     []parentOrderParam `json:"parameters"`
}

//ParentOrderAcceptanceID Parent order acceptance ID
type ParentOrderAcceptanceID string

type parentOrderResponseBody struct {
	ID ParentOrderAcceptanceID `json:"parent_order_acceptance_id"`
}

func newParentOrderParam(productCode string, conditionType string, side string, price float32, triggerPrice float32, size float32) *parentOrderParam {
	return &parentOrderParam{productCode, conditionType, side, price, triggerPrice, size}
}

func parentOrder(orderMethod string, minuteToExpire int, timeInForce string, parameters []parentOrderParam) ParentOrderAcceptanceID {
	postBody := parentOrderBody{orderMethod, minuteToExpire, timeInForce, parameters}
	client := service.NewBitClient()
	pathDir := "/v1/me/sendparentorder"
	response := new(parentOrderResponseBody)
	client.Post(pathDir, postBody, response)
	return response.ID
}

//IFDStopMarketOrder Do parent order of IFD, STOP order first and then MARKET order of same size
func IFDStopMarketOrder(side string, triggerPrice float32, size float32) ParentOrderAcceptanceID {
	productType := "FX_BTC_JPY"
	firstOrder := newParentOrderParam(productType, "STOP", side, 0, triggerPrice, size)
	secondOrder := newParentOrderParam(productType, "MARKET", side, 0, 0, size)
	parameters := []parentOrderParam{*firstOrder, *secondOrder}
	return parentOrder("IFD", 2, "GTC", parameters)
}

//StopOrder Do parent order of STOP
func StopOrder(side string, triggerPrice float32, size float32) ParentOrderAcceptanceID {
	order := newParentOrderParam("FX_BTC_JPY", "STOP", side, 0, triggerPrice, size)
	return parentOrder("SIMPLE", 2, "GTC", []parentOrderParam{*order})
}
