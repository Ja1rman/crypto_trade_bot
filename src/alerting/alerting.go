package alerting

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	
	"crypto_trading/src/logger"
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
