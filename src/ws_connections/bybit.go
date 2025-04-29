package ws_connections

import (
	"fmt"
	"encoding/json"

	bybit "github.com/bybit-exchange/bybit.go.api"

	"crypto_trading/src/logger"
	"crypto_trading/src/handlers"
)

type OrderBookJsonMessage struct {
	Topic string            `json:"topic"`
	Type  string            `json:"type"`
	Ts    int64             `json:"ts"`
	Data  handlers.OrderBookJsonData `json:"data"`
	Cts   int64             `json:"cts"`
}

func StartByBitConnection(tickers []string) {
	ws := bybit.NewBybitPublicWebSocket("wss://stream.bybit.com/v5/public/spot", func(message string) error {
		//fmt.Println("Received message: ", message)
		var jsonMessage OrderBookJsonMessage
		if err := json.Unmarshal([]byte(message), &jsonMessage); err != nil {
			logger.Logger.Println("Ошибка парсинга JSON:", err)
			return err
		}
		// для задержек
		//eventTime := time.UnixMilli(jsonMessage.Cts)
		//now := time.Now()
		//diff := now.Sub(eventTime).Milliseconds()
		//if diff > 100 && jsonMessage.Ts != 0 {
			//go alerting.SendMessage(fmt.Sprintf("Задержка принятия сообщения %d ms", diff))
			//fmt.Printf("Задержка принятия сообщения %d ms\n", diff)
		//}
		go handlers.UpdateOrdersBook(jsonMessage.Data, jsonMessage.Ts)
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
	select {}
}
