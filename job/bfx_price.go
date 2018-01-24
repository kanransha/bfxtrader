package job

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fxtrader/model"
	"net/http"
	"time"
)

type bfxBoard struct {
	MidPrice float32    `json:"mid_price"`
	Bids     []bfxOrder `json:"bids"`
	Asks     []bfxOrder `json:"asks"`
}

type bfxOrder struct {
	Price float32 `json:"price"`
	Size  float32 `json:"size"`
}

func bfxBoardRequest() *bytes.Buffer {
	res, err := http.Get("https://api.bitflyer.jp/v1/board?product_code=FX_BTC_JPY")
	if err != nil {
		return bfxBoardRequest()
	}
	if res.StatusCode != 200 {
		fmt.Printf("Board StatusCode = %d\n", res.StatusCode)
		res.Body.Close()
		return bfxBoardRequest()
	}
	defer res.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	return buf
}

//GetBFXPrice Get BFX price and write to channel
func GetBFXPrice(c chan *model.BFXPrice) {
	jsonBytes := bfxBoardRequest().Bytes()
	timeNow := time.Now()
	jsonData := new(bfxBoard)
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
