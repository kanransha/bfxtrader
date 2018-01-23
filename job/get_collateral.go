package job

import (
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
	jsonData := new(CollateralValues)
	client.Get("/v1/me/getcollateral", "", jsonData)
	return jsonData
}
