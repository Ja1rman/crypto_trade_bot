package ws_connections

import (
	"crypto_trading/src/logger"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"

	"crypto_trading/websocket-proto/mexc/protobuf"
	"crypto_trading/src/stats"
	"crypto_trading/src/handlers"
)

var (
	maxMexcTime int64 = 0
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

func connectToMexc(tickers []string, prices *handlers.Prices) error {
	const reconnectDelay = 5 * time.Second

	url := "wss://wbs-api.mexc.com/ws"
	params := make([]string, 0, len(tickers))
	for _, ticker := range tickers {
		params = append(params, fmt.Sprintf("spot@public.bookTicker.batch.v3.api.pb@%s", ticker))
	}

	for {
		logger.Logger.Printf("Подключение к WebSocket для тикеров: %v", tickers)
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			logger.Logger.Println("Ошибка соединения с WebSocket:", err)
			time.Sleep(reconnectDelay)
			continue
		}

		if err := conn.WriteJSON(map[string]interface{}{
			"method": "SUBSCRIPTION",
			"params": params,
		}); err != nil {
			logger.Logger.Println("Ошибка подписки:", err)
			conn.Close()
			time.Sleep(reconnectDelay)
			continue
		}

		pingTicker := time.NewTicker(30 * time.Second)

		readErr := make(chan error, 1)
		go func() {
			for {
				conn.SetReadDeadline(time.Now().Add(1 * time.Second))
				msgType, message, err := conn.ReadMessage()
				if err != nil {
					readErr <- err
					return
				}
				if msgType != websocket.BinaryMessage {
					continue
				}

				var msg protobuf.PushDataV3ApiWrapper
				if err := proto.Unmarshal(message, &msg); err != nil {
					logger.Logger.Println("Ошибка разбора protobuf:", err)
					continue
				}
				eventTime := time.UnixMilli(msg.GetSendTime())
				now := time.Now()
				diff := now.Sub(eventTime).Milliseconds()
				if diff > maxMexcTime && msg.GetSendTime() != 0 {
					maxMexcTime = diff
				}
				if now.Unix() % 15 == 0 {
					go stats.PushToPrometheus(stats.RecieveFromServerDuration, "mexc", float64(maxMexcTime), "duration")
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
		}()

	loop:
		for {
			select {
			case <-pingTicker.C:
				if err := conn.WriteJSON(map[string]interface{}{"method": "PING"}); err != nil {
					logger.Logger.Println("Ошибка при отправке PING:", err)
					break loop
				}
			case err := <-readErr:
				logger.Logger.Println("Ошибка чтения WebSocket:", err)
				break loop
			}
		}

		conn.Close()
		pingTicker.Stop()
		logger.Logger.Println("Переподключение через", reconnectDelay)
		time.Sleep(reconnectDelay)
	}
}
