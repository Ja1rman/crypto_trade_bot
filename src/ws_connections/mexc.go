package ws_connections

import (
	"encoding/json"
	"fmt"
	"crypto_trading/src/logger"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"crypto_trading/websocket-proto/mexc/protobuf"

	"crypto_trading/src/handlers"
)


func StartMexcConnection(tickers []string) {
	url := "wss://wbs-api.mexc.com/ws"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		logger.Logger.Fatal("Ошибка соединения с WebSocket:", err)
	}
	defer conn.Close()

	logger.Logger.Println("WebSocket подключён")

	const batchSize = 30
	for i := 0; i < len(tickers); i += batchSize {
		end := i + batchSize
		if end > len(tickers) {
			end = len(tickers)
		}
		params := make([]string, 0)
		for _, ticker := range tickers[i:end] {
			params = append(params, fmt.Sprintf("spot@public.bookTicker.batch.v3.api.pb@%s", ticker))
		}

		sub := map[string]interface{}{
			"method": "SUBSCRIPTION",
			"params": params,
		}

		if err := conn.WriteJSON(sub); err != nil {
			logger.Logger.Println("Ошибка подписки:", err)
			continue
		}
	}

	go func() {
		for {
			msgType, message, err := conn.ReadMessage()
			if err != nil {
				logger.Logger.Println("Ошибка при чтении WebSocket:", err)
			}
			if msgType != websocket.BinaryMessage {
				continue
			}
			var msg protobuf.PushDataV3ApiWrapper
			err = proto.Unmarshal(message, &msg)
			if err != nil {
				logger.Logger.Println("Ошибка при разборе protobuf:", err)
			}
			var data handlers.OrderBookJsonData
			data.Symbol = msg.GetSymbol()
			data.Seq = msg.GetSendTime()
			data.Update = data.Seq
			publicBookTickerArr := msg.GetPublicBookTickerBatch().GetItems()
			if len(publicBookTickerArr) == 0 {
				continue
			}
			publicBookTicker := publicBookTickerArr[0]
			data.Bids = [][]string{{publicBookTicker.GetBidprice(), publicBookTicker.GetBidquantity()}}
			data.Asks = [][]string{{publicBookTicker.GetAskprice(), publicBookTicker.GetAskquantity()}}
			go handlers.UpdateOrdersBook(data, data.Seq)
		}
	}()
	go func() {
		pinhTicker := time.NewTicker(30 * time.Second)
		defer pinhTicker.Stop()
	
		for range pinhTicker.C {
			ping := map[string]string{
				"method": "PING",
			}
			pingData, _ := json.Marshal(ping)
	
			err := conn.WriteMessage(websocket.TextMessage, pingData)
			if err != nil {
				logger.Logger.Println("Ошибка при отправке PING:", err)
				return
			}
			logger.Logger.Println("PING отправлен")
		}
	}()
	
	select {}
}
