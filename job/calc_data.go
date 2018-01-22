package job

import (
	"sync"

	"../model"
)

//CalcBFXData Calculate BFX Data
func CalcBFXData(c chan *model.BFXPrice, m *sync.Mutex, market *model.BFXMarket) {
	currentPrice := <-c
	m.Lock()
	market.SetBFXPrice(currentPrice)
	m.Unlock()
}
