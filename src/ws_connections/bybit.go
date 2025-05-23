package ws_connections

import (
	"fmt"
	"time"
	"encoding/json"

	bybit "github.com/bybit-exchange/bybit.go.api"

	"crypto_trading/src/logger"
	"crypto_trading/src/handlers"
	"crypto_trading/src/stats"
)

type OrderBookJsonMessage struct {
	Topic string            `json:"topic"`
	Type  string            `json:"type"`
	Ts    int64             `json:"ts"`
	Data  handlers.OrderBookJsonData `json:"data"`
	Cts   int64             `json:"cts"`
}

var (
	maxBybitTime int64 = 0
)

func StartByBitConnection(tickers []string, prices *handlers.Prices) {
	ws := bybit.NewBybitPublicWebSocket("wss://stream.bybit.com/v5/public/spot", func(message string) error {
		var jsonMessage OrderBookJsonMessage
		if err := json.Unmarshal([]byte(message), &jsonMessage); err != nil {
			logger.Logger.Println("Ошибка парсинга JSON:", err)
			return err
		}
		eventTime := time.UnixMilli(jsonMessage.Ts)
		diff := time.Since(eventTime).Milliseconds()
		if diff > maxBybitTime && jsonMessage.Ts != 0 {
			maxBybitTime = diff
		}
		go prices.UpdateOrdersBook(jsonMessage.Data, jsonMessage.Ts)
		return nil
	})
	_ = ws.Connect()

	// Разбиение подписок на группы по 10
	const batchSize = 10
	for i := 0; i < len(tickers); i += batchSize {
		end := i + batchSize
		if end > len(tickers) {
			end = len(tickers)
		}
		subscriptions := make([]string, 0, end-i)
		for _, t := range tickers[i:end] {
			subscriptions = append(subscriptions, fmt.Sprintf("orderbook.1.%s", t))
		}
		_, _ = ws.SendSubscription(subscriptions)
	}
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			stats.PushToPrometheus(stats.RecieveFromServerDuration, "bybit", float64(maxBybitTime), "duration")
			maxBybitTime = 0
		}
	}()
	select {}
}
