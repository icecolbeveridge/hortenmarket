package hortenmarket

import (
	"fmt"
	"math/rand"
	"time"
)

type Candle struct { // this is a common object and it probably wants to be defined somewhere sensible.
	High float64
	Low float64
	Open float64
	Close float64
	Volume float64
	// Time time.Time
}

func randomPrice(value float64) float64 {
	e := 1. + rand.NormFloat64() * 0.05
	return value * e
}	

func mins(arr []float64) (float64, int) {
	m := 0.
	mi := -1
	for i, e := range(arr) {
		if i == 0 || e < m {
			m = e
			mi = i
		}
	}
	return m, mi
}

func maxs(arr []float64) (float64, int) {
	m := 0.
	mi := -1
	for i, e := range(arr) {
		if i == 0 || e > m {
			m = e
			mi = i
		}
	}
	return m, mi
}

func remove(elem int, arr []float64) []float64 {
	arr[elem] = arr[len(arr)-1]
	return arr[:len(arr)-1]
}

// candlestickMaker simulates a market 
func candlestickMaker(ouch chan interface{}) {
	offers := make([]float64, 0)
	bids := make([]float64, 0)
	trades := make([]float64, 0)
	newCandle := true
	price := 100. // starting value

	Open := price
	High := price
	Close := price
	Low := price
	Volume := 0.

	for i := 0; ; i ++   {
		if i % 100 == 0 {
			// new candle
			newCandle = true			
		}
		newBid := randomPrice(price)
		newOffer := randomPrice(price)
		offers = append(offers, newOffer)
		bids = append(bids, newBid)
		
		minOffer, moi := mins(offers)
		maxBid, mbi := maxs(bids)

		if maxBid > minOffer {
			price = (maxBid+minOffer)*0.5
			trades = append(trades, price)
			if newCandle {
				Open = price
				High = price
				Close = price
				Low = price
				newCandle = false
				Volume = 100.
			} else {
				Close = price
				if price > High {High = price}
				if price < Low {Low = price}
				Volume += 100.
			}
			offers = remove(moi, offers)
			bids = remove(mbi, bids)
		}
		if i % 100 == 99 {
			// emit a candle!
			<- time.After(1000 * time.Millisecond)
			ouch <- Candle{Open:Open, Close:Close, High:High, Low:Low, Volume:Volume, Time: time.Now()}
			newCandle = true
			
		}

	}

}

func listen(ch chan Candle) {
	for c := range(ch) {
		fmt.Println("Open", c.Open)
		fmt.Println("High", c.High)
		fmt.Println("Low", c.Low)
		fmt.Println("Close", c.Close)
		fmt.Println("Volume", c.Volume)

		
	}
}

func Run() {
	ch := make(chan Candle)
	go listen(ch)
	candlestickMaker(ch)
}

// call from elsewhere as hortenmarket.Run()
