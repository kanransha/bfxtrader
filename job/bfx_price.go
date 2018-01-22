package job

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"../model"
)

type board struct {
	MidPrice float32 `json:"mid_price"`
	Bids     []order `json:"bids"`
	Asks     []order `json:"asks"`
}

type order struct {
	Price float32 `json:"price"`
	Size  float32 `json:"size"`
}

func boardRequest() *bytes.Buffer {
	res, err := http.Get("https://api.bitflyer.jp/v1/board?product_code=FX_BTC_JPY")
	if err != nil {
		return boardRequest()
	}
	if res.StatusCode != 200 {
		fmt.Printf("Board StatusCode = %d\n", res.StatusCode)
		res.Body.Close()
		return boardRequest()
	}
	defer res.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	return buf
}

//GetBFXPrice Get BFX price and write to channel
func GetBFXPrice(c chan *model.BFXPrice) {
	jsonBytes := boardRequest().Bytes()
	timeNow := time.Now()
	jsonData := new(board)
	if err := json.Unmarshal(jsonBytes, jsonData); err != nil {
		fmt.Println("JSON Unmarchal error")
		GetBFXPrice(c)
		return
	}
	if jsonData.MidPrice == 0 {
		fmt.Println("JSON MidPrice 0")
		GetBFXPrice(c)
		return
	}
	c <- model.NewBFXPrice(timeNow, jsonData.MidPrice)
}
