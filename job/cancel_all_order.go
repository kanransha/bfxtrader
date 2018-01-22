package job

import (
	"encoding/json"
	"fmt"

	"../service"
)

type cancelAllOrderBody struct {
	ProductCode string `json:"product_code"`
}

//CancelAllOrder Cancel all order
func CancelAllOrder() {
	jsonBytes, err := json.Marshal(cancelAllOrderBody{"FX_BTC_JPY"})
	client := service.NewBitClient()
	request, err := client.NewRequest("/v1/me/cancelallchildorders", "POST", string(jsonBytes))
	if err != nil {
		fmt.Println("CancelAllOrder Request Error")
		CancelAllOrder()
		return
	}
	res, err := client.Do(request)
	if err != nil {
		fmt.Println("CancelAllOrder Response Error")
		CancelAllOrder()
		return
	}
	if res.StatusCode != 200 {
		res.Body.Close()
		fmt.Printf("CancelAllOrder StatusCode = %d\n", res.StatusCode)
		CancelAllOrder()
		return
	}
	res.Body.Close()
}
