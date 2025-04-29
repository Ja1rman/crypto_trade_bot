package utils

import (
	"crypto_trading/src/config"
)

func GenerateCurrencyRoutes(infoMap map[string]config.Info) map[string]map[string]struct{} {
	routes := make(map[string]map[string]struct{})

	for _, info := range infoMap {
		base := info.BaseCoin
		quote := info.QuoteCoin
		if _, exists := routes[base]; !exists {
			routes[base] = make(map[string]struct{})
		}
		routes[base][quote] = struct{}{}
		if _, exists := routes[quote]; !exists {
			routes[quote] = make(map[string]struct{})
		}
		routes[quote][base] = struct{}{}
	}

	return filterRoutes(routes)
}

func filterRoutes(routes map[string]map[string]struct{}) map[string]map[string]struct{} {
	filtered := make(map[string]map[string]struct{})

	for base, quotes := range routes {
		if len(quotes) >= 2 {
			filtered[base] = quotes
		}
	}

	for base, quotes := range filtered {
		filteredQuotes := make(map[string]struct{})
		for quote := range quotes {
			if _, exists := filtered[quote]; exists {
				filteredQuotes[quote] = struct{}{}
			}
		}
		filtered[base] = filteredQuotes
	}

	return filtered
}
