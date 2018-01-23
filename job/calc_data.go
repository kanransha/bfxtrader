package job

import (
	"bfxtrader/model"
	"sync"
)

//CalcBFXData Calculate BFX Data
func CalcBFXData(c chan *model.BFXPrice, finish chan bool, m *sync.Mutex, market *model.BFXMarket) {
	currentPrice := <-c
	m.Lock()
	market.SetBFXPrice(currentPrice)
	m.Unlock()
	finish <- true
}
