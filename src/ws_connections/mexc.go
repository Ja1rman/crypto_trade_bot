package ws_connections

import (
	"crypto_trading/src/logger"
	"fmt"
	"time"

	"crypto_trading/websocket-proto/mexc/protobuf"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"

	"crypto_trading/src/handlers"
)

func StartMexcConnection(tickers []string, prices *handlers.Prices) {
	const batchSize = 10
	for i := 0; i < len(tickers); i += batchSize {
		end := i + batchSize
		if end > len(tickers) {
			end = len(tickers)
		}
		batch := tickers[i:end]
		go connectToMexc(batch, prices)
	}
	select {}
}

func connectToMexc(tickers []string, prices *handlers.Prices) {
	url := "wss://wbs-api.mexc.com/ws"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		logger.Logger.Fatal("Ошибка соединения с WebSocket:", err)
		return
	}
	defer conn.Close()

	params := make([]string, 0)
	for _, ticker := range tickers {
		params = append(params, fmt.Sprintf("spot@public.bookTicker.batch.v3.api.pb@%s", ticker))
	}
	sub := map[string]interface{}{
		"method": "SUBSCRIPTION",
		"params": params,
	}
	if err := conn.WriteJSON(sub); err != nil {
		logger.Logger.Println("Ошибка подписки:", err)
		return
	}
	pingTicker := time.NewTicker(30 * time.Second)
	defer pingTicker.Stop()

	for {
		select {
		case <-pingTicker.C:
			ping := map[string]interface{}{
				"method": "PING",
			}
			if err := conn.WriteJSON(ping); err != nil {
				logger.Logger.Println("Ошибка при отправке PING:", err)
				return
			}
		default:
			msgType, message, err := conn.ReadMessage()
			if err != nil {
				logger.Logger.Println("Ошибка при чтении WebSocket:", err)
				return
			}
			if msgType != websocket.BinaryMessage {
				continue
			}
			var msg protobuf.PushDataV3ApiWrapper
			err = proto.Unmarshal(message, &msg)
			if err != nil {
				logger.Logger.Println("Ошибка при разборе protobuf:", err)
				continue
			}
			publicBookTickerArr := msg.GetPublicBookTickerBatch().GetItems()
			if len(publicBookTickerArr) == 0 {
				continue
			}
			publicBookTicker := publicBookTickerArr[0]
			data := handlers.OrderBookJsonData{
				Symbol: msg.GetSymbol(),
				Seq:    msg.GetSendTime(),
				Update: msg.GetSendTime(),
				Bids:   [][]string{{publicBookTicker.GetBidprice(), publicBookTicker.GetBidquantity()}},
				Asks:   [][]string{{publicBookTicker.GetAskprice(), publicBookTicker.GetAskquantity()}},
			}
			prices.UpdateOrdersBook(data, data.Seq)
		}
	}
}
