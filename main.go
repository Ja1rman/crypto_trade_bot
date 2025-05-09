package main

import (
	"time"	

	"crypto_trading/src/config"
	"crypto_trading/src/logger"
	"crypto_trading/src/utils"
	"crypto_trading/src/ws_connections"
)

func generateTickersList(tickers map[string]config.Info) []string {
	var result []string
	for ticker := range tickers {
		result = append(result, ticker)
	}
	return result
}

func main() {
	logger.LoggerSetup()
	//mexcCurrencySettings := utils.SetMexcSettings()
	bybitCurrencySettings := utils.SetBybitSettings()

	//tickers := generateTickersList(mexcCurrencySettings.SUBSCRIBE_TICKERS_LIST)
	tickers := generateTickersList(bybitCurrencySettings.SUBSCRIBE_TICKERS_LIST)

	//go ws_connections.StartMexcConnection(tickers, &mexcCurrencySettings.PRICES)
	go ws_connections.StartByBitConnection(tickers, &bybitCurrencySettings.PRICES)
	time.Sleep(40 * time.Second)
	//go mexcCurrencySettings.LaunchInfiniteAnalyze()
	go bybitCurrencySettings.LaunchInfiniteAnalyze()
	select {}
}
