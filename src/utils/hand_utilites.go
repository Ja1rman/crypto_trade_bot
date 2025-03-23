package utils

// all tickers curl -L -X GET 'https://api.bybit.com/v5/market/tickers?category=spot'
import (
	"fmt"
    "context"
	"os"
	"encoding/json"
	"strconv"

    bybit "github.com/wuhewuhe/bybit.go.api"

	"crypto_trading/src/logger"
	"crypto_trading/src/config"
)

type InstrumentInfo struct {
	Symbol      string `json:"symbol"`
	BaseCoin    string `json:"baseCoin"`
	QuoteCoin   string `json:"quoteCoin"`
	LotSizeFilter struct {
		BasePrecision  string `json:"basePrecision"`
		QuotePrecision string `json:"quotePrecision"`
	} `json:"lotSizeFilter"`
}

type InstrumentInfoResult struct {
	Category string           `json:"category"`
	List     []InstrumentInfo `json:"list"`
}


func writeToFile(values interface{}, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		logger.Logger.Fatal(err)
	}
	defer file.Close()

	data, err := json.MarshalIndent(values, "", "  ")
	if err != nil {
		logger.Logger.Fatal(err)
	}

	_, err = file.Write(data)
	if err != nil {
		logger.Logger.Fatal(err)
	}
}

func getTickers(apiKey string, apiSecret string) {
	client := bybit.NewBybitHttpClient(apiKey, apiSecret, bybit.WithBaseURL(bybit.MAINNET))
	params := map[string]interface{}{"category": "spot"}
	results, _ := client.NewUtaBybitServiceWithParams(params).GetMarketTickers(context.Background())
	writeToFile(results, "results.json")
}

func GetCurrenciesPrecision() {
	client := bybit.NewBybitHttpClient("", "", bybit.WithBaseURL(bybit.MAINNET))
	precisionMap := make(map[string]config.Info)

	for symbol := range config.SUBSCRIBE_TICKERS_LIST {
		params := map[string]interface{}{"category": "spot", "symbol": symbol}
		response, err := client.NewUtaBybitServiceWithParams(params).GetInstrumentInfo(context.Background())
		if err != nil {
			fmt.Printf("Ошибка при получении информации для символа %s: %v\n", symbol, err)
			continue
		}

		if response.RetCode != 0 {
			fmt.Printf("API вернуло ошибку для символа %s: %s\n", symbol, response.RetMsg)
			continue
		}

		var result InstrumentInfoResult
		resultBytes, err := json.Marshal(response.Result)
		if err != nil {
			fmt.Printf("Ошибка при преобразовании result в JSON: %v\n", err)
			continue
		}
		err = json.Unmarshal(resultBytes, &result)
		if err != nil {
			fmt.Printf("Ошибка при разборе result: %v\n", err)
			continue
		}

		for _, instrument := range result.List {
			basePrecision, _ := strconv.ParseFloat(instrument.LotSizeFilter.BasePrecision, 64)
			quotePrecision, _ := strconv.ParseFloat(instrument.LotSizeFilter.QuotePrecision, 64)

			precisionMap[symbol] = config.Info{
				BaseCoin:  instrument.BaseCoin,
				QuoteCoin: instrument.QuoteCoin,
				Precision: config.Precision{
					BasePrecision: basePrecision,
					QuotePrecision: quotePrecision,
				},
			}
		}
	}

	// Запись в файл
	writeToFile(precisionMap, "results.json")
}
