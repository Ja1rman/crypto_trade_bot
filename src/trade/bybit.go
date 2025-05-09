package trade

import (
	"fmt"
	"context"
	"strconv"
	"encoding/json"

	bybit "github.com/bybit-exchange/bybit.go.api"

	"crypto_trading/src/config"
	"crypto_trading/src/logger"
)

var (
	client = bybit.NewBybitHttpClient(config.API_KEY, config.API_KEY_SECRET, bybit.WithBaseURL(config.NET))
)


func CreateBybitOrder(symbol string, price float64, qty float64, coinType string, side string, orderType string) (string, error) {
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

	logInfo := fmt.Sprintf("params: %v", params)
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

func GetBybitOrderInfo(orderId string) ([]map[string]interface{}, error) {
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

func GetBybitWalletBalance(coin string) (map[string]float64, error) {
	params := map[string]interface{}{"accountType": "UNIFIED", "coin": coin}
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

func SumBybitQty(orders []map[string]interface{}, inversion bool) float64 {
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
		cumExecFeeStr, ok := order["cumExecFee"].(string)
		if !ok {
			logger.Logger.Println("Ошибка: cumExecFee не является строкой")
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
		cumExecFee, err := strconv.ParseFloat(cumExecFeeStr, 64)
		if err != nil {
			logger.Logger.Println("Ошибка конвертации cumExecFee:", err)
			continue
		}
		if inversion {
			res += qty * avgPrice
		} else {
			res += qty
		}
		res -= cumExecFee
	}
	return res
}
