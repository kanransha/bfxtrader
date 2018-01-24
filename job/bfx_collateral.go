package job

import (
	"fxtrader/service"
)

//BFXCollateralValues Collateral values of bitflyer FX
type BFXCollateralValues struct {
	Collateral        float32 `json:"collateral"`
	OpenPositionPNL   float32 `json:"open_position_pnl"`
	RequireCollateral float32 `json:"require_collateral"`
	KeepRate          float32 `json:"keep_rate"`
}

//GetBFXCollateralValues Get collateral values by bitflyer API
func GetBFXCollateralValues() *BFXCollateralValues {
	client := service.NewBitClient()
	jsonData := new(BFXCollateralValues)
	client.Get("/v1/me/getcollateral", "", jsonData)
	return jsonData
}
