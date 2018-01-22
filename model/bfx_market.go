package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"time"
)

type ohlcData []float64

type ohlcResponse struct {
	Result    map[string]([]ohlcData) `json:"result"`
	Allowance struct {
		Cost      int64 `json:"cost"`
		Remaining int64 `json:"remaining"`
	} `json:"allowance"`
}

//BFXMarket Array of BFXValues of certain period
type BFXMarket struct {
	dataCount     int
	slowEMAPeriod int
	fastEMAPeriod int
	secInterval   int
	dataArray     []BFXValues
}

func (ohlc *ohlcData) Time() time.Time {
	return time.Unix(int64((*ohlc)[0]), 0)
}

func (ohlc *ohlcData) EndPrice() float32 {
	return float32((*ohlc)[4])
}

func ohlcRequest(count int, interval int) *bytes.Buffer {
	rawURL := "https://api.cryptowat.ch/markets/bitflyer/btcfxjpy/ohlc"
	timeBefore := time.Now().Unix()
	timeAfter := timeBefore - int64(interval)*(int64(count)+1)
	params := url.Values{}
	params.Add("before", fmt.Sprint(timeBefore))
	params.Add("after", fmt.Sprint(timeAfter))
	params.Add("periods", fmt.Sprint(interval))
	req, errReq := http.NewRequest("GET", rawURL, nil)
	if errReq != nil {
		fmt.Println("Error in OHLC Request New Request")
		return ohlcRequest(count, interval)
	}
	req.URL.RawQuery = params.Encode()
	client := &http.Client{}
	res, errRes := client.Do(req)
	if errRes != nil {
		fmt.Println("Error in OHLC Request Response")
		return ohlcRequest(count, interval)
	}
	if res.StatusCode != 200 {
		fmt.Printf("OHLC StatusCode = %d\n", res.StatusCode)
		res.Body.Close()
		fmt.Println(req.URL.String())
		panic("OHLC")
		//return ohlcRequest(count, interval)
	}
	defer res.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	return buf
}

func getInitialPrices(count int, interval int) []BFXPrice {
	jsonBytes := ohlcRequest(count, interval).Bytes()
	jsonData := new(ohlcResponse)
	if err := json.Unmarshal(jsonBytes, jsonData); err != nil {
		fmt.Println("OHLC JSON Unmarchal error")
		panic("OHLC")
	}
	priceArray := BFXPrices{}
	for _, ohlc := range jsonData.Result[fmt.Sprint(interval)] {
		priceArray = append(priceArray, BFXPrice{ohlc.Time(), ohlc.EndPrice()})
		if ohlc.EndPrice() == 0 {
			panic("Wrong Price Value")
		}
	}
	sort.Sort(priceArray)
	priceArray = priceArray[len(priceArray)-count:]
	return priceArray
}

func calcInitialEMA(count int, priceArray []BFXPrice, period int) []float32 {
	emaArray := make([]float32, count)
	firstSum := float32(0)
	for i := 0; i < period; i++ {
		firstSum += priceArray[i].value
		emaArray[i] = 0
	}
	emaArray[period-1] = firstSum / float32(period)
	for i := period; i < count; i++ {
		emaArray[i] = calcNextEMA(priceArray[i].value, emaArray[i-1], period)
	}
	return emaArray
}

//NewBFXMarket Create New BFXMarket
func NewBFXMarket(count int, slowEMAPeriod int, fastEMAPeriod int, secInterval int) *BFXMarket {
	if slowEMAPeriod < fastEMAPeriod {
		fmt.Println("Warning (slowEmaPeriod < fastEMAPeriod): EMA Period Values are Swapped")
		tempEMAPeriod := slowEMAPeriod
		slowEMAPeriod = fastEMAPeriod
		fastEMAPeriod = tempEMAPeriod
	}
	dataArray := make([]BFXValues, count)
	priceArray := getInitialPrices(count, secInterval)
	slowEMAArray := calcInitialEMA(count, priceArray, slowEMAPeriod)
	fastEMAArray := calcInitialEMA(count, priceArray, fastEMAPeriod)
	for i := 0; i < count; i++ {
		dataArray[i].price = priceArray[i]
		dataArray[i].slowEMA = slowEMAArray[i]
		dataArray[i].fastEMA = fastEMAArray[i]
		dataArray[i].Print()
	}
	return &BFXMarket{count, slowEMAPeriod, fastEMAPeriod, secInterval, dataArray}
}

//GetLastValues Get last BFXValues in BFXMarket
func (market *BFXMarket) GetLastValues() *BFXValues {
	count := market.dataCount
	return &(market.dataArray[count-1])
}

//SecInterval Get interval as second
func (market *BFXMarket) SecInterval() int {
	return market.secInterval
}

//SetBFXPrice Set BFXPrice and calculate BFXValues
func (market *BFXMarket) SetBFXPrice(price *BFXPrice) {
	lastValues := market.GetLastValues()
	newSlowEMA := calcNextEMA(price.value, lastValues.slowEMA, market.slowEMAPeriod)
	newFastEMA := calcNextEMA(price.value, lastValues.fastEMA, market.fastEMAPeriod)
	newBFXValues := BFXValues{*price, newSlowEMA, newFastEMA}
	newBFXValues.Print()
	market.dataArray = market.dataArray[1:]
	market.dataArray = append(market.dataArray, newBFXValues)
}
