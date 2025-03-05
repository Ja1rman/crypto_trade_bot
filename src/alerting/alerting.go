package alerting

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	
	"crypto_trading/src/logger"
	"crypto_trading/src/handlers"
)

var (
	TELEGRAM_TOKEN = os.Getenv("TELEGRAM_TOKEN")
	CHAT_ID = "-1001646996694"
)

func SendMessage(message string) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", TELEGRAM_TOKEN)

	data := url.Values{}
	data.Set("chat_id", CHAT_ID)
	data.Set("text", message)
	data.Set("parse_mode", "Markdown") // Можно использовать HTML или Markdown

	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		logger.Logger.Printf("Ошибка при отправке сообщения: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Logger.Printf("Неудачная попытка отправки сообщения, статус: %s", resp.Status)
		return
	}
}


func PrintPrices() {
	handlers.PRICES.Lock()
	for key := range handlers.PRICES.Cache {
		logger.Logger.Printf("Key: %s\n", key)
		for subKey, orderBookData := range handlers.PRICES.Cache[key] {
			logger.Logger.Printf("  SubKey: %s\n", subKey)
			logger.Logger.Printf("    Ask: Price=%.10f, Size=%.10f, Seq=%d\n",
				orderBookData.Ask.Price, orderBookData.Ask.Size, orderBookData.Ask.Seq)
				logger.Logger.Printf("    Bid: Price=%.10f, Size=%.10f, Seq=%d\n",
				orderBookData.Bid.Price, orderBookData.Bid.Size, orderBookData.Bid.Seq)
		}
	}
	handlers.PRICES.Unlock()
}
