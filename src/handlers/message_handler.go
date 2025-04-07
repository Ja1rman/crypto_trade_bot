package handlers

import (
	"strconv"
	"sync"

	"crypto_trading/src/logger"
	"crypto_trading/src/config"
)

type Prices struct {
	sync.Mutex
	Cache map[string]map[string]OrderBookData
}

var (
	PRICES Prices = Prices{
		Cache: make(map[string]map[string]OrderBookData),
	}
)

type StockExchangeGlassNote struct {
	Price float64
	Size float64
	Seq int64
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

func UpdateOrdersBook(msg OrderBookJsonData) {
	key := msg.Symbol
	value, exists := config.SUBSCRIBE_TICKERS_LIST[key]
	if !exists {
		return
	}

	hasBids := len(msg.Bids) > 0
	hasAsks := len(msg.Asks) > 0

	if !hasBids && !hasAsks {
		logger.Logger.Println("Нет данных в Bids и Asks для обновления.")
		return
	}

	newItem := OrderBookData{
		Ask: StockExchangeGlassNote{Price: 0, Size: 0, Seq: 0},
		Bid: StockExchangeGlassNote{Price: 0, Size: 0, Seq: 0},
	}

	if hasBids {
		bestBid := msg.Bids[0]
		newItem.Bid = StockExchangeGlassNote{
			Price: ParseFloat(bestBid[0]),
			Size: ParseFloat(bestBid[1]),
			Seq: msg.Seq,
		}
	}

	if hasAsks {
		bestAsk := msg.Asks[0]
		newItem.Ask = StockExchangeGlassNote{
			Price: ParseFloat(bestAsk[0]),
			Size: ParseFloat(bestAsk[1]),
			Seq: msg.Seq,
		}
	}

	updatePriceMap(value.QuoteCoin, value.BaseCoin, newItem)
	newItem2 := OrderBookData{
		Ask: StockExchangeGlassNote{Price: 0, Size: 0, Seq: 0},
		Bid: StockExchangeGlassNote{Price: 0, Size: 0, Seq: 0},
	}
	newItem2.Ask.Size = newItem.Bid.Size * newItem.Bid.Price
	if newItem.Bid.Price > 0 {
		newItem2.Ask.Price = 1 / newItem.Bid.Price
	}
	newItem2.Bid.Size = newItem.Ask.Size * newItem.Ask.Price
	if newItem.Ask.Price > 0 {
		newItem2.Bid.Price = 1 / newItem.Ask.Price
	}
	updatePriceMap(value.BaseCoin, value.QuoteCoin, newItem2)
}

func updatePriceMap(first string, second string, newItem OrderBookData) {
	PRICES.Lock()
	if _, exists := PRICES.Cache[first]; !exists {
		PRICES.Cache[first] = make(map[string]OrderBookData)
	}

	if oldItem, exists := PRICES.Cache[first][second]; !exists {
		PRICES.Cache[first][second] = newItem
	} else {
		if oldItem.Bid.Seq < newItem.Bid.Seq {
			oldItem.Bid = newItem.Bid
		}
		if oldItem.Ask.Seq < newItem.Ask.Seq {
			oldItem.Ask = newItem.Ask
		}
		PRICES.Cache[first][second] = oldItem
	}
	PRICES.Unlock()
}

func DeepCopyCache() map[string]map[string]OrderBookData {
	PRICES.Lock()
	defer PRICES.Unlock()
	res := make(map[string]map[string]OrderBookData)
	for key1, subMap := range PRICES.Cache {
		newSubMap := make(map[string]OrderBookData)
		for key2, orderBook := range subMap {
			newSubMap[key2] = OrderBookData{
				Ask: StockExchangeGlassNote{
					Price: orderBook.Ask.Price,
					Size:  orderBook.Ask.Size,
					Seq:   orderBook.Ask.Seq,
				},
				Bid: StockExchangeGlassNote{
					Price: orderBook.Bid.Price,
					Size:  orderBook.Bid.Size,
					Seq:   orderBook.Bid.Seq,
				},
			}
		}
		res[key1] = newSubMap
	}
	return res
}
