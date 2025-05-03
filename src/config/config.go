package config

import (
	"os"
	
	bybit "github.com/bybit-exchange/bybit.go.api"
)

var (
	API_KEY string = os.Getenv("API_KEY")
	API_KEY_SECRET string = os.Getenv("API_KEY_SECRET")
)

type Info struct {
	BaseCoin string `json:"BaseCoin"`
	QuoteCoin string `json:"QuoteCoin"`
	Precision Precision `json:"Precision"`
	TickSize float64 `json:"TickSize"`
}

type Precision struct {
	BasePrecision  float64 `json:"BasePrecision"`
	QuotePrecision float64 `json:"QuotePrecision"`
}

var (
	NET = bybit.MAINNET
	TRADE bool = false
)
