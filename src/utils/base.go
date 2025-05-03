package utils

import (
	"crypto_trading/src/config"
)

func GenerateRoutesAndInfoMap(infoMap map[string]config.Info) (map[string]map[string]struct{}, map[string]config.Info) {
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
	routes = filterRoutes(routes)
	infoMap = filterInfoMapByRoutes(infoMap, routes)
	return routes, infoMap
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

func filterInfoMapByRoutes(infoMap map[string]config.Info, routes map[string]map[string]struct{}) map[string]config.Info {
	filtered := make(map[string]config.Info)
	for symbol, info := range infoMap {
		if _, exists := routes[info.BaseCoin]; !exists {
			continue
		}
		if _, exists := routes[info.QuoteCoin]; !exists {
			continue
		}
		filtered[symbol] = info
	}
	return filtered
}
