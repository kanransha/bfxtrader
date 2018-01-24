package job

import (
	"fxtrader/service"
)

type cancelAllBFXOrderBody struct {
	ProductCode string `json:"product_code"`
}

//CancelAllBFXOrder Cancel all order
func CancelAllBFXOrder() {
	b := cancelAllBFXOrderBody{"FX_BTC_JPY"}
	c := service.NewBitClient()
	c.Post("/v1/me/cancelallchildorders", b, nil)
}
