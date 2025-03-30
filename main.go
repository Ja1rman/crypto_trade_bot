package main

import (
	"encoding/json"
	"fmt"
	//"time"
	//"context"
	bybit "github.com/wuhewuhe/bybit.go.api"

	"crypto_trading/src/logger"
	"crypto_trading/src/handlers"
	"crypto_trading/src/config"
	//"crypto_trading/src/alerting"
	"crypto_trading/src/analyzer"

	//"crypto_trading/src/trade"
	//"crypto_trading/src/utils"
)


type OrderBookJsonMessage struct {
	Topic string            `json:"topic"`
	Type  string            `json:"type"`
	Ts    int64             `json:"ts"`
	Data  handlers.OrderBookJsonData `json:"data"`
	Cts   int64             `json:"cts"`
}

func generateTickersList() []string {
	var result []string
	for ticker := range config.SUBSCRIBE_TICKERS_LIST {
		result = append(result, fmt.Sprintf("orderbook.1.%s", ticker))
	}
	return result
}

func startConnection() {
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
		go handlers.UpdateOrdersBook(jsonMessage.Data)
		return nil
	})
	_ = ws.Connect()

	// Разбиение подписок на группы по 10
	tickers := generateTickersList()
	const batchSize = 10
	for i := 0; i < len(tickers); i += batchSize {
		end := i + batchSize
		if end > len(tickers) {
			end = len(tickers)
		}
		_, _ = ws.SendSubscription(tickers[i:end])
	}
	select {}
}


func main() {
	logger.LoggerSetup()
	go startConnection()
	go analyzer.StartAnalyzies()
	select {}
}
