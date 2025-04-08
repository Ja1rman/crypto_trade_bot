package analyzer_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	_ "fmt"
	
	"crypto_trading/src/analyzer"
	"crypto_trading/src/handlers"
)

для тестов ProcessAnalyze
currencies := []string{"USDT", "MANA", "BTC"}
a := handlers.OrderBookData{
	Ask: handlers.StockExchangeGlassNote{Price: 0.2157, Size: 13900.14, Seq: 2},
	Bid: handlers.StockExchangeGlassNote{Price: 0.2158, Size: 108.06, Seq: 2},
}
b := handlers.OrderBookData{
	Ask: handlers.StockExchangeGlassNote{Price: 0.00000290, Size: 1280.1, Seq: 2},
	Bid: handlers.StockExchangeGlassNote{Price: 0.00000295, Size: 1280.1, Seq: 2},
}
c := handlers.OrderBookData{
	Ask: handlers.StockExchangeGlassNote{Price: 79680.3, Size: 5.1, Seq: 2},
	Bid: handlers.StockExchangeGlassNote{Price: 79690.4, Size: 4.9, Seq: 2},
}
handlers.UpdatePriceMap("MANAUSDT", a)
handlers.UpdatePriceMap("MANABTC", b)
handlers.UpdatePriceMap("BTCUSDT", c)
analyzer.ProcessAnalyze(currencies)


var _ = Describe("Trading Calculations", func() {
	var balance = 1000.0
	Describe("CalculateMarketSize", func() {
		It("first Bid Size", func() {
			tradingPairs := analyzer.TradingPairs{
				handlers.OrderBookData{
					Ask: handlers.StockExchangeGlassNote{Price: 9, Size: 50, Seq: 1},
					Bid: handlers.StockExchangeGlassNote{Price: 10, Size: 40, Seq: 1},
				},
				handlers.OrderBookData{
					Ask: handlers.StockExchangeGlassNote{Price: 10, Size: 50, Seq: 1},
					Bid: handlers.StockExchangeGlassNote{Price: 2, Size: 40, Seq: 1},
				},
				handlers.OrderBookData{
					Ask: handlers.StockExchangeGlassNote{Price: 10, Size: 50, Seq: 1},
					Bid: handlers.StockExchangeGlassNote{Price: 20, Size: 40, Seq: 1},
				},
			}
			expectedSize := 40.0
			res := analyzer.CalculateMarketSize(tradingPairs, balance)
			Expect(res).To(Equal(expectedSize))
		})

		It("BALANCE / first.Bid.Price", func() {
			tradingPairs := analyzer.TradingPairs{
				handlers.OrderBookData{
					Ask: handlers.StockExchangeGlassNote{Price: 9, Size: 50, Seq: 1},
					Bid: handlers.StockExchangeGlassNote{Price: 10, Size: 400, Seq: 1},
				},
				handlers.OrderBookData{
					Ask: handlers.StockExchangeGlassNote{Price: 10, Size: 50, Seq: 1},
					Bid: handlers.StockExchangeGlassNote{Price: 2, Size: 4000, Seq: 1},
				},
				handlers.OrderBookData{
					Ask: handlers.StockExchangeGlassNote{Price: 10, Size: 500, Seq: 1},
					Bid: handlers.StockExchangeGlassNote{Price: 20, Size: 40, Seq: 1},
				},
			}
			expectedSize := 100.0
			res := analyzer.CalculateMarketSize(tradingPairs, balance)
			Expect(res).To(Equal(expectedSize))
		})

		It("By Second Bid Size", func() {
			tradingPairs := analyzer.TradingPairs{
				handlers.OrderBookData{
					Ask: handlers.StockExchangeGlassNote{Price: 9, Size: 50, Seq: 1},
					Bid: handlers.StockExchangeGlassNote{Price: 10, Size: 400, Seq: 1},
				},
				handlers.OrderBookData{
					Ask: handlers.StockExchangeGlassNote{Price: 10, Size: 50, Seq: 1},
					Bid: handlers.StockExchangeGlassNote{Price: 2, Size: 40, Seq: 1},
				},
				handlers.OrderBookData{
					Ask: handlers.StockExchangeGlassNote{Price: 10, Size: 50, Seq: 1},
					Bid: handlers.StockExchangeGlassNote{Price: 20, Size: 40, Seq: 1},
				},
			}
			expectedSize := 50.0
			res := analyzer.CalculateMarketSize(tradingPairs, balance)
			Expect(res).To(Equal(expectedSize))
		})

		It("Third Bid Size", func() {
			tradingPairs := analyzer.TradingPairs{
				handlers.OrderBookData{
					Ask: handlers.StockExchangeGlassNote{Price: 9, Size: 50, Seq: 1},
					Bid: handlers.StockExchangeGlassNote{Price: 10, Size: 400, Seq: 1},
				},
				handlers.OrderBookData{
					Ask: handlers.StockExchangeGlassNote{Price: 10, Size: 500, Seq: 1},
					Bid: handlers.StockExchangeGlassNote{Price: 2, Size: 40, Seq: 1},
				},
				handlers.OrderBookData{
					Ask: handlers.StockExchangeGlassNote{Price: 10, Size: 50, Seq: 1},
					Bid: handlers.StockExchangeGlassNote{Price: 20, Size: 40, Seq: 1},
				},
			}
			expectedSize := 50.0
			res := analyzer.CalculateMarketSize(tradingPairs, balance)
			Expect(res).To(Equal(expectedSize))
		})
	})

	Describe("AnalizeOfferProfit", func() {
		It("Test correctly data", func() {
			tradingPairs := analyzer.TradingPairs{
				handlers.OrderBookData{
					Ask: handlers.StockExchangeGlassNote{Price: 10, Size: 50, Seq: 0},
					Bid: handlers.StockExchangeGlassNote{Price: 10, Size: 40, Seq: 0},
				},
				handlers.OrderBookData{
					Ask: handlers.StockExchangeGlassNote{Price: 10, Size: 50, Seq: 0},
					Bid: handlers.StockExchangeGlassNote{Price: 1, Size: 40, Seq: 0},
				},
				handlers.OrderBookData{
					Ask: handlers.StockExchangeGlassNote{Price: 0.1, Size: 50, Seq: 0},
					Bid: handlers.StockExchangeGlassNote{Price: 2, Size: 40, Seq: 0},
				},
			}
			startSize := 100.0
			expectedProfit := 0.0
			res := analyzer.AnalizeOfferProfit(tradingPairs, startSize)
			Expect(res).To(Equal(expectedProfit))
		})

		It("Test correctly data", func() {
			tradingPairs := analyzer.TradingPairs{
				handlers.OrderBookData{
					Ask: handlers.StockExchangeGlassNote{Price: 10, Size: 50, Seq: 0},
					Bid: handlers.StockExchangeGlassNote{Price: 0, Size: 40, Seq: 0},
				},
				handlers.OrderBookData{
					Ask: handlers.StockExchangeGlassNote{Price: 10, Size: 50, Seq: 0},
					Bid: handlers.StockExchangeGlassNote{Price: 1, Size: 40, Seq: 0},
				},
				handlers.OrderBookData{
					Ask: handlers.StockExchangeGlassNote{Price: 0.1, Size: 50, Seq: 0},
					Bid: handlers.StockExchangeGlassNote{Price: 2, Size: 40, Seq: 0},
				},
			}
			startSize := 100.0
			expectedProfit := 0.0
			res := analyzer.AnalizeOfferProfit(tradingPairs, startSize)
			Expect(res).To(Equal(expectedProfit))
		})
	})
})
