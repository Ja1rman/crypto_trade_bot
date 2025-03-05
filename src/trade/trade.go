package trade


import (
    "fmt"
	"context"

	bybit "github.com/wuhewuhe/bybit.go.api"

	"crypto_trading/src/config"
)

func CreateOrder(symbol string, price float64, qty float64) {
	client := bybit.NewBybitHttpClient(config.API_KEY, config.API_KEY_SECRET, bybit.WithBaseURL(bybit.TESTNET))
	params := map[string]interface{}{
		"category": "spot", 
		"symbol": symbol,
		"side": "Buy", 
		"orderType": "Limit", 
		"timeInForce": "IOC",
		"qty": qty,
		"marketUnit": "baseCoin", // quoteCoin
		"price": price,
	}
	orderResult, err := client.NewUtaBybitServiceWithParams(params).PlaceOrder(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(bybit.PrettyPrint(orderResult))
}





