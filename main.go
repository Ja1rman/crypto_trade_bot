package main

import (
	"time"

	"crypto_trading/src/config"
	"crypto_trading/src/logger"
	"crypto_trading/src/utils"
	"crypto_trading/src/ws_connections"
)

func generateTickersList(tickers map[string]config.Info) []string {
	result := make([]string, 0, len(tickers))
	for ticker := range tickers {
		result = append(result, ticker)
	}
	return result
}

func main() {
	logger.LoggerSetup()

	mexcCurrencySettings := utils.SetMexcSettings()
	bybitCurrencySettings := utils.SetBybitSettings()

	mexcTickers := generateTickersList(mexcCurrencySettings.SUBSCRIBE_TICKERS_LIST)
	bybitTickers := generateTickersList(bybitCurrencySettings.SUBSCRIBE_TICKERS_LIST)

	go ws_connections.StartMexcConnection(mexcTickers, &mexcCurrencySettings.PRICES)
	go ws_connections.StartByBitConnection(bybitTickers, &bybitCurrencySettings.PRICES)

	time.Sleep(40 * time.Second)

	go mexcCurrencySettings.LaunchInfiniteAnalyze()
	go bybitCurrencySettings.LaunchInfiniteAnalyze()

	select {}
}
