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

var (
	PRICES Prices = Prices{
		Cache: make(map[string]OrderBookData),
	}
)

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

func UpdateOrdersBook(msg OrderBookJsonData, ts int64) {
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
	UpdatePriceMap(symbol, newItem, msg.Update)
}

func UpdatePriceMap(symbol string, newItem OrderBookData, update int64) {
	PRICES.Lock()
	defer PRICES.Unlock()
	if oldItem, exists := PRICES.Cache[symbol]; !exists || update == 1 {
		PRICES.Cache[symbol] = newItem
	} else {
		if oldItem.Bid.Seq < newItem.Bid.Seq {
			oldItem.Bid = newItem.Bid
		}
		if oldItem.Ask.Seq < newItem.Ask.Seq {
			oldItem.Ask = newItem.Ask
		}
		PRICES.Cache[symbol] = oldItem
	}
}

func GetOrderBookData(symbol string) OrderBookData {
	PRICES.RLock()
	defer PRICES.RUnlock()
	if item, exists := PRICES.Cache[symbol]; exists {
		return item
	}
	return OrderBookData{}
}
