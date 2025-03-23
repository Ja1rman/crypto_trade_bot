package trade

import (
	"context"
	"fmt"
	"math"
	"strconv"

	bybit "github.com/wuhewuhe/bybit.go.api"

	"crypto_trading/src/config"
)

var (
	client = bybit.NewBybitHttpClient(config.API_KEY, config.API_KEY_SECRET, bybit.WithBaseURL(config.NET))
)

func SellAll(pairsNames []string, qty float64) error {
	errMsg := fmt.Sprintf("Не удалось продать все валюты, pairsNames: %v", pairsNames)
	for i, pairName := range pairsNames {
		if i == 0 || i == len(pairsNames)-1 {
			continue
		}
		symbol, info, err := findSymbol(pairsNames[i], pairsNames[i+1])
		if err != nil {
			return fmt.Errorf("%s ошибка: %s", errMsg, err)
		}
		qty, err = createSellOrder(info, qty, symbol, pairName)
		if err != nil {
			return fmt.Errorf("%s ошибка: %s", errMsg, err)
		}
	}
	return nil
}

func createSellOrder(info config.Info, size float64, symbol string, coin string) (float64, error) {
	coinType := "baseCoin"
	if info.BaseCoin == coin {
		size = RoundCustomStep(size, info.Precision.BasePrecision)
	} else {
		coinType = "quoteCoin"
		size = RoundCustomStep(size, info.Precision.QuotePrecision)
	}
	orderID, err := SellMarketPrice(symbol, size, coinType)
	if err != nil {
		fmt.Println(err)
		return 0.0, err
	}
	orderInfo, err := GetOrderInfo(orderID)
	if err != nil {
		fmt.Println(err)
		return 0.0, err
	}
	if len(orderInfo) == 0 {
		return 0.0, fmt.Errorf("orderInfo not found")
	}
	if coinType == "baseCoin" {
		return SumQty(orderInfo, false), nil
	} else {
		return SumQty(orderInfo, true), nil
	}
}

func SellMarketPrice(symbol string, qty float64, coinType string) (string, error) {
	params := map[string]interface{}{
		"category":   "spot",
		"symbol":     symbol,
		"side":       "Sell",
		"orderType":  "Market",
		"qty":        fmt.Sprint(qty),
		"marketUnit": coinType, // quoteCoin, baseCoin
	}
	logInfo := fmt.Sprintf("Qty: %f, Symbol: %s", qty, symbol)
	serverResponse, err := client.NewUtaBybitServiceWithParams(params).PlaceOrder(context.Background())
	if err != nil || serverResponse.RetCode != 0 {
		return "", fmt.Errorf("invalid request. RetCode: %d, RetMsg: %s, LogInfo: %s, Error: %v", serverResponse.RetCode, serverResponse.RetMsg, logInfo, err)
	}
	orderResult, ok := serverResponse.Result.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("orderResult not found in response. LogInfo: %s", logInfo)
	}
	orderID, ok := orderResult["orderId"].(string)
	if !ok {
		return "", fmt.Errorf("orderId not found in response. LogInfo: %s", logInfo)
	}
	return orderID, nil
}

func ProcessFirstPair(size float64, price float64, pairsNames []string) (float64, error) {
	symbol, info, err := findSymbol(pairsNames[0], pairsNames[1])
	if err != nil {
		return 0.0, err
	}
	return createFirstOrder(info, size, price, symbol, pairsNames[1])
}

func findSymbol(firstCurrency string, secondCurrency string) (string, config.Info, error) {
	symbol := firstCurrency + secondCurrency
	if info, exists := config.SUBSCRIBE_TICKERS_LIST[symbol]; exists {
		return symbol, info, nil
	}
	symbol = secondCurrency + firstCurrency
	if info, exists := config.SUBSCRIBE_TICKERS_LIST[symbol]; exists {
		return symbol, info, nil
	}
	return "", config.Info{}, fmt.Errorf("Pairs %s not found in config", symbol)
}

func createFirstOrder(info config.Info, size float64, price float64, symbol string, coin string) (float64, error) {
	side := "Buy"
	if info.BaseCoin == coin {
		side = "Sell"
		size = size * price
		price = 1.0 / price
		size = RoundCustomStep(size, info.Precision.QuotePrecision)
	} else {
		size = RoundCustomStep(size, info.Precision.BasePrecision)
	}
	//price = RoundCustomStep(price, info.Precision.QuotePrecision)
	orderID, err := CreateLimitOrder(symbol, price, size, side)
	if err != nil {
		fmt.Println(err)
		return 0.0, err
	}
	orderInfo, err := GetOrderInfo(orderID)
	if err != nil {
		fmt.Println(err)
		return 0.0, err
	}
	if len(orderInfo) == 0 {
		return 0.0, fmt.Errorf("orderInfo not found")
	}
	if side == "Buy" {
		return SumQty(orderInfo, false), nil
	} else {
		return SumQty(orderInfo, true), nil
	}
}

func CreateLimitOrder(symbol string, price float64, qty float64, side string) (string, error) {

	params := map[string]interface{}{
		"category":    "spot",
		"symbol":      symbol,
		"side":        side,
		"orderType":   "Limit",
		"timeInForce": "IOC",
		"qty":         fmt.Sprint(qty),
		"marketUnit":  "baseCoin", // quoteCoin, baseCoin
		"price":       fmt.Sprint(price),
	}
	logInfo := fmt.Sprintf("Price: %f, Qty: %f, Symbol: %s", price, qty, symbol)
	serverResponse, err := client.NewUtaBybitServiceWithParams(params).PlaceOrder(context.Background())
	if err != nil || serverResponse.RetCode != 0 {
		return "", fmt.Errorf("invalid request. RetCode: %d, RetMsg: %s, LogInfo: %s, Error: %v", serverResponse.RetCode, serverResponse.RetMsg, logInfo, err)
	}
	orderResult, ok := serverResponse.Result.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("orderResult not found in response. LogInfo: %s", logInfo)
	}
	orderID, ok := orderResult["orderId"].(string)
	if !ok {
		return "", fmt.Errorf("orderId not found in response. LogInfo: %s", logInfo)
	}
	return orderID, nil
}

func RoundCustomStep(number, step float64) float64 {
	return math.Floor(number/step) * step
}

func SumQty(orders []map[string]interface{}, invertion bool) float64 {
	res := 0.
	for _, order := range orders {
		qtyStr, ok := order["cumExecQty"].(string)
		if !ok {
			fmt.Println("Ошибка: qty не является строкой")
			continue
		}
		avgPriceStr, ok := order["avgPrice"].(string)
		if !ok {
			fmt.Println("Ошибка: avgPrice не является строкой")
			continue
		}
		qty, err := strconv.ParseFloat(qtyStr, 64)
		if err != nil {
			fmt.Println("Ошибка конвертации qty:", err)
			continue
		}
		avgPrice, err := strconv.ParseFloat(avgPriceStr, 64)
		if err != nil {
			fmt.Println("Ошибка конвертации qty:", err)
			continue
		}
		if invertion {
			res += qty * avgPrice
		} else {
			res += qty
		}
	}
	return res
}

func GetOrderInfo(orderId string) ([]map[string]interface{}, error) {
	params := map[string]interface{}{
		"category": "spot",
		"orderId":  orderId,
	}
	serverResponse, err := client.NewUtaBybitServiceWithParams(params).GetOpenOrders(context.Background())
	if err != nil || serverResponse.RetCode != 0 {
		return []map[string]interface{}{}, fmt.Errorf("invalid request. RetCode: %d, RetMsg: %s, orderId: %s, Error: %v", serverResponse.RetCode, serverResponse.RetMsg, orderId, err)
	}
	orderResult, ok := serverResponse.Result.(map[string]interface{})
	if !ok {
		return []map[string]interface{}{}, fmt.Errorf("orderResult not found in response. orderId: %s", orderId)
	}
	list, ok := orderResult["list"].([]interface{})
	if !ok {
		return []map[string]interface{}{}, fmt.Errorf("Ошибка: orderResult[\"list\"] не является []interface{}. orderId: %s", orderId)
	}

	// Преобразуем каждый элемент списка в map[string]interface{}
	orderInfo := make([]map[string]interface{}, 0)
	for _, item := range list {
		if order, ok := item.(map[string]interface{}); ok {
			orderInfo = append(orderInfo, order)
		} else {
			return []map[string]interface{}{}, fmt.Errorf("Ошибка: элемент в list не является map[string]interface{}. orderId: %s", orderId)
		}
	}
	return orderInfo, nil
}
