package handlers

import (
	"strconv"
	"sync"

	"crypto_trading/src/logger"
)

type Prices struct {
	sync.RWMutex
	Cache map[string]OrderBookData
}

type StockExchangeGlassNote struct {
	Price float64
	Size float64
	Seq int64
	Ts int64
}

type OrderBookData struct {
	Ask StockExchangeGlassNote
	Bid StockExchangeGlassNote
}

type OrderBookJsonData struct {
	Symbol string `json:"s"`
	Bids [][]string `json:"b"`
	Asks [][]string `json:"a"`
	Update int64 `json:"u"`
	Seq int64 `json:"seq"`
}

func ParseFloat(s string) float64 {
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		logger.Logger.Println("Error parsing float:", err)
	}
	return value
}

func (prices *Prices) UpdateOrdersBook(msg OrderBookJsonData, ts int64) {
	symbol := msg.Symbol

	hasBids := len(msg.Bids) > 0
	hasAsks := len(msg.Asks) > 0

	if !hasBids && !hasAsks {
		return
	}

	newItem := OrderBookData{
		Ask: StockExchangeGlassNote{Price: 0, Size: 0, Seq: 0, Ts: ts},
		Bid: StockExchangeGlassNote{Price: 0, Size: 0, Seq: 0, Ts: ts},
	}

	if hasBids {
		bestBid := msg.Bids[0]
		newItem.Bid = StockExchangeGlassNote{
			Price: ParseFloat(bestBid[0]),
			Size: ParseFloat(bestBid[1]),
			Seq: msg.Seq,
			Ts: ts,
		}
	}

	if hasAsks {
		bestAsk := msg.Asks[0]
		newItem.Ask = StockExchangeGlassNote{
			Price: ParseFloat(bestAsk[0]),
			Size: ParseFloat(bestAsk[1]),
			Seq: msg.Seq,
			Ts: ts,
		}
	}
	prices.UpdatePriceMap(symbol, newItem, msg.Update)
}

func (prices *Prices) UpdatePriceMap(symbol string, newItem OrderBookData, update int64) {
	prices.Lock()
	defer prices.Unlock()
	if oldItem, exists := prices.Cache[symbol]; !exists || update == 1 {
		prices.Cache[symbol] = newItem
	} else {
		if oldItem.Bid.Seq < newItem.Bid.Seq {
			oldItem.Bid = newItem.Bid
		}
		if oldItem.Ask.Seq < newItem.Ask.Seq {
			oldItem.Ask = newItem.Ask
		}
		prices.Cache[symbol] = oldItem
	}
}

func (prices *Prices) GetOrderBookData(symbol string) OrderBookData {
	prices.RLock()
	defer prices.RUnlock()
	if item, exists := prices.Cache[symbol]; exists {
		return item
	}
	return OrderBookData{}
}

func (prices *Prices) GetLenOfPrices() int {
	prices.RLock()
	defer prices.RUnlock()
	return len(prices.Cache)
}
