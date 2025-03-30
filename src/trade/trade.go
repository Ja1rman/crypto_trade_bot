package trade

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"encoding/json"

	bybit "github.com/wuhewuhe/bybit.go.api"

	"crypto_trading/src/config"
	"crypto_trading/src/logger"
)

var (
	client = bybit.NewBybitHttpClient(config.API_KEY, config.API_KEY_SECRET, bybit.WithBaseURL(config.NET))
)

func ProcessCycle(startSize float64, firstOrderPrice float64, pairsNames []string) {
	qty, err := ProcessFirstPair(startSize, firstOrderPrice, pairsNames)
	if err!= nil || qty <= 0.0 {
		logger.Logger.Println(err)
		return
	}
	pairsNames = append(pairsNames, pairsNames[0])
	SellAll(pairsNames, qty)
}

func SellAll(pairsNames []string, qty float64) error {
	errMsg := fmt.Sprintf("Не удалось продать все валюты, pairsNames: %v", pairsNames)
	for i, pairName := range pairsNames {
		if i == 0 || i == len(pairsNames)-1 {
			continue
		}
		symbol, info, err := FindSymbol(pairsNames[i], pairsNames[i+1])
		if err != nil {
			return fmt.Errorf("%s ошибка: %s, %v", errMsg, err, info)
		}
		qty, err = CreateSellOrder(info, qty, symbol, pairName)
		if err != nil {
			return fmt.Errorf("%s ошибка: %s", errMsg, err)
		}
	}
	return nil
}

func CreateSellOrder(info config.Info, size float64, symbol string, coin string) (float64, error) {
	coinType := "baseCoin"
	side := "Sell"
	if info.BaseCoin == coin {
		size = RoundCustomStep(size, info.Precision.BasePrecision)
	} else {
		coinType = "quoteCoin"
		side = "Buy"
		size = RoundCustomStep(size, info.Precision.QuotePrecision)
	}
	orderID, err := CreateOrder(symbol, 1., size, coinType, side, "Market")
	if err != nil {
		logger.Logger.Println(err)
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
	if coinType == "quoteCoin" {
		return SumQty(orderInfo, false), nil
	} else {
		return SumQty(orderInfo, true), nil
	}
}


func ProcessFirstPair(size float64, price float64, pairsNames []string) (float64, error) {
	symbol, info, err := FindSymbol(pairsNames[0], pairsNames[1])
	if err != nil {
		return 0.0, err
	}
	return createFirstOrder(info, size, price, symbol, pairsNames[1])
}

func FindSymbol(firstCurrency string, secondCurrency string) (string, config.Info, error) {
	symbol := firstCurrency + secondCurrency
	if info, exists := config.SUBSCRIBE_TICKERS_LIST[symbol]; exists {
		return symbol, info, nil
	}
	symbol = secondCurrency + firstCurrency
	if info, exists := config.SUBSCRIBE_TICKERS_LIST[symbol]; exists {
		return symbol, info, nil
	}
	return "", config.Info{}, fmt.Errorf("pairs %s not found in config", symbol)
}

func createFirstOrder(info config.Info, size float64, price float64, symbol string, coin string) (float64, error) {
	side := "Buy"
	coinType := "baseCoin"
	if info.BaseCoin != coin {
		side = "Sell"
		size = size * price
		price = 1.0 / price
		size = RoundCustomStep(size, info.Precision.QuotePrecision)
	} else {
		size = RoundCustomStep(size, info.Precision.BasePrecision)
	}
	price = RoundCustomStep(price, info.TickSize)

	orderID, err := CreateOrder(symbol, price, size, coinType, side, "Limit")
	if err != nil {
		logger.Logger.Println(err)
		return 0.0, err
	}
	orderInfo, err := GetOrderInfo(orderID)
	if err != nil {
		logger.Logger.Println(err)
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

func CreateOrder(symbol string, price float64, qty float64, coinType string, side string, orderType string) (string, error) {

	params := map[string]interface{}{
		"category":    "spot",
		"symbol":      symbol,
		"side":        side,
		"orderType":   orderType, // Market, Limit
		"qty":         fmt.Sprint(qty),
		"marketUnit":  coinType, // quoteCoin, baseCoin
		
	}
	if orderType == "Limit" {
		params["price"] = fmt.Sprint(price)
		params["timeInForce"] = "IOC"
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
	return math.Round(number*step) / step
}

func SumQty(orders []map[string]interface{}, invertion bool) float64 {
	res := 0.
	for _, order := range orders {
		qtyStr, ok := order["cumExecQty"].(string)
		if !ok {
			logger.Logger.Println("Ошибка: qty не является строкой")
			continue
		}
		avgPriceStr, ok := order["avgPrice"].(string)
		if !ok {
			logger.Logger.Println("Ошибка: avgPrice не является строкой")
			continue
		}
		qty, err := strconv.ParseFloat(qtyStr, 64)
		if err != nil {
			logger.Logger.Println("Ошибка конвертации qty:", err)
			continue
		}
		avgPrice, err := strconv.ParseFloat(avgPriceStr, 64)
		if err != nil {
			logger.Logger.Println("Ошибка конвертации qty:", err)
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
		return []map[string]interface{}{}, fmt.Errorf("ошибка: orderResult[\"list\"] не является []interface{}. orderId: %s", orderId)
	}

	// Преобразуем каждый элемент списка в map[string]interface{}
	orderInfo := make([]map[string]interface{}, 0)
	for _, item := range list {
		if order, ok := item.(map[string]interface{}); ok {
			orderInfo = append(orderInfo, order)
		} else {
			return []map[string]interface{}{}, fmt.Errorf("ошибка: элемент в list не является map[string]interface{}. orderId: %s", orderId)
		}
	}
	return orderInfo, nil
}

func SaveBookInfo(symbol string) {
	client := bybit.NewBybitHttpClient("", "", bybit.WithBaseURL(config.NET))
	params := map[string]interface{}{"category": "spot", "symbol": symbol}
	response, err := client.NewUtaBybitServiceWithParams(params).GetOrderBookInfo(context.Background())
	if err != nil {
		logger.Logger.Printf("Ошибка при получении информации для символа %s: %v\n", symbol, err)
		return
	}
	logger.Logger.Println(response)
}

func GetWalletBalance() (map[string]float64, error) {
	params := map[string]interface{}{"accountType": "UNIFIED", "coin": "USDT,USDC"}
	accountResult, err := client.NewUtaBybitServiceWithParams(params).GetAccountWallet(context.Background())
	if err != nil {
		return map[string]float64{}, fmt.Errorf("ошибка получения баланса: %s", err)
	}

	jsonData, err := json.Marshal(accountResult)
	if err != nil {
		return map[string]float64{}, fmt.Errorf("ошибка маршалинга: %s", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(jsonData, &response); err != nil {
		return map[string]float64{}, fmt.Errorf("ошибка парсинга JSON: %s", err)
	}

	result, ok := response["result"].(map[string]interface{})
	if !ok {
		return map[string]float64{}, fmt.Errorf("нет result в ответе")
	}

	list, ok := result["list"].([]interface{})
	if !ok || len(list) == 0 {
		return map[string]float64{}, fmt.Errorf("нет list в ответе")
	}

	balanceMap := make(map[string]float64)

	for _, item := range list {
		account, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		coins, ok := account["coin"].([]interface{})
		if !ok {
			continue
		}

		for _, c := range coins {
			coinData, ok := c.(map[string]interface{})
			if !ok {
				continue
			}

			coinName, _ := coinData["coin"].(string)
			walletBalanceStr, _ := coinData["walletBalance"].(string)

			// Конвертация строки в float64
			walletBalance, err := strconv.ParseFloat(walletBalanceStr, 64)
			if err != nil {
				logger.Logger.Println("Ошибка конвертации баланса для", coinName, ":", err)
				continue
			}

			balanceMap[coinName] = walletBalance
		}
	}

	return balanceMap, nil
}
