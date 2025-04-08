package utils

// all tickers curl -L -X GET 'https://api.bybit.com/v5/market/tickers?category=spot'
import (
	"fmt"
    "context"
	"os"
	"encoding/json"
	"strconv"
	"strings"

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
	PriceFilter struct {
		TickSize string `json:"tickSize"`
	} `json:"priceFilter"`
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

func GetCurrenciesPrecision() {  // Сделать обратную точность (1/x)
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
			tickSize, _ := strconv.ParseFloat(instrument.PriceFilter.TickSize, 64)
			if tickSize != 0.01 {
				fmt.Println(tickSize)
			}
			precisionMap[symbol] = config.Info{
				BaseCoin:  instrument.BaseCoin,
				QuoteCoin: instrument.QuoteCoin,
				Precision: config.Precision{
					BasePrecision: basePrecision,
					QuotePrecision: quotePrecision,
				},
				TickSize: tickSize,
			}
		}
	}

	// Запись в файл
	writeToFile(precisionMap, "results.json")
}

func GenerateRoutes() {
	// Генерируем карту валютных пар
	currencyPairs := GenerateCurrencyPairs()

	// Создаем содержимое файла
	content := `var CURRENCY_PAIRS = map[string][]string{`

	// Добавляем пары в содержимое файла
	for currency, pairs := range currencyPairs {
		content += fmt.Sprintf("\t\"%s\": {%s},\n", currency, formatPairs(pairs))
	}

	content += "}\n"

	// Записываем в файл
	err := os.WriteFile("currency_pairs.go", []byte(content), 0644)
	if err != nil {
		fmt.Printf("Ошибка записи файла: %v\n", err)
		return
	}

	fmt.Println("Файл currency_pairs.go успешно создан")
}

func formatPairs(pairs []string) string {
	quoted := make([]string, len(pairs))
	for i, pair := range pairs {
		quoted[i] = fmt.Sprintf("\"%s\"", pair)
	}
	return strings.Join(quoted, ", ")
}

func GenerateCurrencyPairs() map[string][]string {
	result := make(map[string][]string)
    
	for _, info := range config.SUBSCRIBE_TICKERS_LIST {
		if _, exists := result[info.BaseCoin]; !exists {
			result[info.BaseCoin] = []string{}
		}
		result[info.BaseCoin] = append(result[info.BaseCoin], info.QuoteCoin)
        
		if _, exists := result[info.QuoteCoin]; !exists {
			result[info.QuoteCoin] = []string{}
		}
		result[info.QuoteCoin] = append(result[info.QuoteCoin], info.BaseCoin)
	}
    
	for currency, pairs := range result {
		uniquePairs := make(map[string]bool)
		for _, pair := range pairs {
			uniquePairs[pair] = true
		}
        
		result[currency] = make([]string, 0, len(uniquePairs))
		for pair := range uniquePairs {
			result[currency] = append(result[currency], pair)
		}
	}
    
	return result
}

