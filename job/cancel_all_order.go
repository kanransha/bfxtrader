package job

import (
	"../service"
)

type cancelAllOrderBody struct {
	ProductCode string `json:"product_code"`
}

//CancelAllOrder Cancel all order
func CancelAllOrder() {
	postBody := cancelAllOrderBody{"FX_BTC_JPY"}
	client := service.NewBitClient()
	client.Post("/v1/me/cancelallchildorders", postBody, nil)
}
