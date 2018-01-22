package job

import (
	"bytes"
	"encoding/json"
	"fmt"

	"../service"
)

//CollateralValues Collateral values of bitflyer FX
type CollateralValues struct {
	Collateral        float32 `json:"collateral"`
	OpenPositionPNL   float32 `json:"open_position_pnl"`
	RequireCollateral float32 `json:"require_collateral"`
	KeepRate          float32 `json:"keep_rate"`
}

//GetCollateralValues Get collateral values by bitflyer API
func GetCollateralValues() *CollateralValues {
	client := service.NewBitClient()
	request, err := client.NewRequest("/v1/me/getcollateral", "GET", "")
	if err != nil {
		fmt.Println("GetCollateralValues Request Error")
		return GetCollateralValues()
	}
	res, err := client.Do(request)
	if err != nil {
		fmt.Println("GetCollateralValues Response Error")
		return GetCollateralValues()
	}
	if res.StatusCode != 200 {
		fmt.Printf("Collateral StatusCode = %d\n", res.StatusCode)
		res.Body.Close()
		return GetCollateralValues()
	}
	defer res.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	jsonBytes := buf.Bytes()
	jsonData := new(CollateralValues)
	if err := json.Unmarshal(jsonBytes, jsonData); err != nil {
		fmt.Println("Collateral JSON Unmarchal error")
		return GetCollateralValues()
	}
	return jsonData
}
