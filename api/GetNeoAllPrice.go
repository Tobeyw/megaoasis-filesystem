package api

import (
	"metaoasis-filesystem/cache"
)

func GetAllPrice(asset string) (cache.Price, error) {
	USD, err := GetPriceFromCoinMarketMap("USD", asset)
	if err != nil {
		return cache.Price{}, err
	}
	KRW, err := GetPriceFromCoinMarketMap("KRW", asset)
	if err != nil {
		return cache.Price{}, err
	}
	IDR, err := GetPriceFromCoinMarketMap("IDR", asset)
	if err != nil {
		return cache.Price{}, err
	}
	THB, err := GetPriceFromCoinMarketMap("THB", asset)
	if err != nil {
		return cache.Price{}, err
	}
	SGD, err := GetPriceFromCoinMarketMap("SGD", asset)
	if err != nil {
		return cache.Price{}, err
	}

	return cache.Price{
		USD: USD,
		KRW: KRW,
		IDR: IDR,
		THB: THB,
		SGD: SGD,
	}, nil
}
