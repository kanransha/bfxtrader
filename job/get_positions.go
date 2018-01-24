package job

import (
	"bfxtrader/service"
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
	pathDir := "/v1/me/getpositions"
	queryStr := "product_code=FX_BTC_JPY"
	jsonData := new(Positions)
	client.Get(pathDir, queryStr, jsonData)
	return jsonData
}

//GetCurrentBFX Get current BFX size and side
func GetCurrentBFX() (float32, string) {
	positions := GetPositions()
	if len(*positions) == 0 {
		return float32(0), "ZERO"
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
