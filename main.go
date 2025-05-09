package main

import (
	"time"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"

	"crypto_trading/src/config"
	//"crypto_trading/src/logger"
	//"crypto_trading/src/utils"
	//"crypto_trading/src/ws_connections"
)

func generateTickersList(tickers map[string]config.Info) []string {
	var result []string
	for ticker := range tickers {
		result = append(result, ticker)
	}
	return result
}

// func main() {
// 	logger.LoggerSetup()
// 	mexcCurrencySettings := utils.SetMexcSettings()
// 	//bybitCurrencySettings := utils.SetBybitSettings()

// 	tickers := generateTickersList(mexcCurrencySettings.SUBSCRIBE_TICKERS_LIST)

// 	go ws_connections.StartMexcConnection(tickers, &mexcCurrencySettings.PRICES)
// 	//go ws_connections.StartByBitConnection(tickers, &mexcCurrencySettings.PRICES)
// 	time.Sleep(40 * time.Second)
// 	go mexcCurrencySettings.LaunchInfiniteAnalyze()
// 	//go bybitCurrencySettings.LaunchInfiniteAnalyze()
// 	select {}
// }

var functionDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
	Name:    "function_duration_seconds",
	Help:    "Время выполнения функции в секундах",
	Buckets: prometheus.DefBuckets,
})

func main() {
	for {
		start := time.Now()
		time.Sleep(500 * time.Millisecond) // Имитация работы
		duration := time.Since(start).Seconds()
		functionDuration.Observe(duration)

		if err := push.New("http://pushgateway:9091", "my_job").
			Collector(functionDuration).
			Grouping("instance", "go-app").
			Push(); err != nil {
			fmt.Println("Could not push metrics:", err)
		}

		time.Sleep(2 * time.Second)
	}
}
