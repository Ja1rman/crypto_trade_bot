package trade

import (
	"fmt"
	"strconv"
	"strings"
	"encoding/json"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"time"
	
	"crypto_trading/src/config"
	"crypto_trading/src/logger"
)

type OrderResponse struct {
	OrderID string `json:"orderId"`
}

type OrderInfoResponse struct {
	Symbol              string `json:"symbol"`
	OrderID             int64  `json:"orderId"`
	OrderListID         int64  `json:"orderListId"`
	ClientOrderID       string `json:"clientOrderId"`
	Price               string `json:"price"`
	Qty                 string `json:"Qty"`
	ExecutedQty         string `json:"executedQty"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	Status              string `json:"status"`
	TimeInForce         string `json:"timeInForce"`
	Type                string `json:"type"`
	Side                string `json:"side"`
	StopPrice           string `json:"stopPrice"`
	Time                int64  `json:"time"`
	UpdateTime          int64  `json:"updateTime"`
	IsWorking           bool   `json:"isWorking"`
	OrigQuoteOrderQty   string `json:"origQuoteOrderQty"`
}

type Balance struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

type AccountResponse struct {
	CanTrade    bool      `json:"canTrade"`
	CanWithdraw bool      `json:"canWithdraw"`
	CanDeposit  bool      `json:"canDeposit"`
	UpdateTime  *int64    `json:"updateTime"`
	AccountType string    `json:"accountType"`
	Balances    []Balance `json:"balances"`
	Permissions []string  `json:"permissions"`
}


func GetMexcOrderInfo(symbol string, orderId string) (*OrderInfoResponse, error) {
	endpoint := "/api/v3/order"
	params := url.Values{}
	params.Add("orderId", orderId)
	params.Add("symbol", symbol)
	timestamp := time.Now().UnixMilli()
	params.Add("timestamp", strconv.FormatInt(timestamp, 10))
	signature := signParams(params.Encode(), config.SECRET_KEY_MEXC)
	params.Add("signature", signature)

	req, err := http.NewRequest("GET", config.MEXC_URL+endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("X-MEXC-APIKEY", config.ACCESS_KEY_MEXC)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status code: %d\n%s", resp.StatusCode, body)
	}

	var order OrderInfoResponse
	if err := json.Unmarshal(body, &order); err != nil {
		return nil, err
	}

	return &order, nil

}

func GetMexcWalletBalance() (map[string]float64, error) {
	endpoint := "/api/v3/account"

	params := url.Values{}
	timestamp := time.Now().UnixMilli()
	params.Add("timestamp", strconv.FormatInt(timestamp, 10))
	signature := signParams(params.Encode(), config.SECRET_KEY_MEXC)
	params.Add("signature", signature)

	req, err := http.NewRequest("GET", config.MEXC_URL+endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-MEXC-APIKEY", config.ACCESS_KEY_MEXC)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status code: %d\n%s", resp.StatusCode, body)
	}
	var account AccountResponse
	if err := json.Unmarshal(body, &account); err != nil {
		return nil, err
	}

	result := make(map[string]float64)
	for _, b := range account.Balances {
		free, err := strconv.ParseFloat(b.Free, 64)
		if err != nil {
			continue // пропускаем, если не число
		}
		result[b.Asset] = free
	}

	return result, nil
}

func CreateMexcOrder(symbol string, price float64, quantity float64, coinType string, side string, orderType string) (string, error) {
	endpoint := "/api/v3/order"

	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("side", strings.ToUpper(side))
	if orderType == "Limit" {
		params.Add("type", "IMMEDIATE_OR_CANCEL")
	} else {
		params.Add("type", strings.ToUpper(orderType))
	}

	if coinType == "baseCoin" {
		params.Add("quantity", fmt.Sprint(quantity))
	}
	if coinType == "quoteCoin" {
		params.Add("quoteOrderQty", fmt.Sprint(quantity))
	}
	if price > 0 {
		params.Add("price", fmt.Sprint(price))
	}

	timestamp := time.Now().UnixMilli()
	params.Add("timestamp", strconv.FormatInt(timestamp, 10))

	signature := signParams(params.Encode(), config.SECRET_KEY_MEXC)
	params.Add("signature", signature)

	req, err := http.NewRequest("POST", config.MEXC_URL+endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("X-MEXC-APIKEY", config.ACCESS_KEY_MEXC)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-200 status code: %d\n%s", resp.StatusCode, body)
	}

	var orderResp OrderResponse
	if err := json.Unmarshal(body, &orderResp); err != nil {
		logger.Logger.Fatalf("failed to parse response: %v", err)
		return "", err
	}

	return orderResp.OrderID, nil
}

func signParams(data, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func SumMexcQty(order *OrderInfoResponse, inversion bool) float64 {
	if order.Status != "FILLED" && order.Status != "PARTIALLY_FILLED" && order.Status != "PARTIALLY_CANCELED" {
		return 0.
	}
	executedQty, err := strconv.ParseFloat(order.ExecutedQty, 64)
	if err != nil {
		logger.Logger.Println("Ошибка конвертации executedQty:", err)
		return 0.
	}
	cummulativeQuoteQty, err := strconv.ParseFloat(order.CummulativeQuoteQty, 64)
	if err != nil {
		logger.Logger.Println("Ошибка конвертации cummulativeQuoteQty:", err)
		return 0.
	}
	if inversion {
		return executedQty 
	}
	return cummulativeQuoteQty
}
