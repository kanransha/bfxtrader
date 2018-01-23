package job

import (
	"bytes"
	"encoding/json"
	"fmt"

	"../service"
)

//Position Position of bitflyer FX
type Position struct {
	ProductCode         string  `json:"product_code"`
	Side                string  `json:"side"`
	Price               float32 `json:"price"`
	Size                float32 `json:"size"`
	Commission          float32 `json:"commission"`
	SwapPointAccumulate float32 `json:"swap_point_accumulate"`
	RequireCollateral   float32 `json:"require_collateral"`
	OpenDate            string  `json:"open_date"`
	Leverage            float32 `json:"leverage"`
	PNL                 float32 `json:"pnl"`
}

//Positions Array of position
type Positions []Position

//GetPositions Get positions
func GetPositions() *Positions {
	client := service.NewBitClient()
	request, err := client.NewRequest("/v1/me/getpositions?product_code=FX_BTC_JPY", "GET", "")
	if err != nil {
		fmt.Println("GetPositions Request Error")
		return GetPositions()
	}
	res, err := client.Do(request)
	if err != nil {
		fmt.Println("GetPositions Response Error")
		return GetPositions()
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	if res.StatusCode != 200 {
		fmt.Printf("GetPositions StatusCode = %d\n", res.StatusCode)
		fmt.Println(buf.String())
		res.Body.Close()
		return GetPositions()
	}
	defer res.Body.Close()
	jsonBytes := buf.Bytes()
	jsonData := new(Positions)
	if err := json.Unmarshal(jsonBytes, jsonData); err != nil {
		fmt.Println("GetPositions JSON Unmarchal error")
		fmt.Println(buf.String())
		return GetPositions()
	}
	return jsonData
}

//GetCurrentBFX Get current BFX size and side
func GetCurrentBFX() (float32, string) {
	positions := GetPositions()
	if len(*positions) == 0 {
		return float32(0), ""
	}
	side := (*positions)[0].Side
	size := float32(0)
	for _, pos := range *positions {
		if pos.Side == side {
			size += pos.Size
		} else {
			size -= pos.Size
		}
	}
	return size, side
}
