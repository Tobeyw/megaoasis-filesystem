package api

import (
	"metaoasis-filesystem/cache"
)

func GetAllPrice(asset string) (cache.CurrencyList, error) {
	USD, usdUpdateTime, err := GetPriceFromCoinMarketMap("USD", asset)
	if err != nil {
		return cache.CurrencyList{}, err
	}
	KRW, krwUpdateTime, err := GetPriceFromCoinMarketMap("KRW", asset)
	if err != nil {
		return cache.CurrencyList{}, err
	}
	IDR, idrUpdateTime, err := GetPriceFromCoinMarketMap("IDR", asset)
	if err != nil {
		return cache.CurrencyList{}, err
	}
	THB, thbUpdateTime, err := GetPriceFromCoinMarketMap("THB", asset)
	if err != nil {
		return cache.CurrencyList{}, err
	}
	SGD, sgdUpdateTime, err := GetPriceFromCoinMarketMap("SGD", asset)
	if err != nil {
		return cache.CurrencyList{}, err
	}

	return cache.CurrencyList{
		USD: cache.Currency{
			Price:      USD,
			UpdateTime: usdUpdateTime,
		},
		KRW: cache.Currency{
			Price:      KRW,
			UpdateTime: krwUpdateTime,
		},
		IDR: cache.Currency{
			Price:      IDR,
			UpdateTime: idrUpdateTime,
		},
		THB: cache.Currency{
			Price:      THB,
			UpdateTime: thbUpdateTime,
		},
		SGD: cache.Currency{
			Price:      SGD,
			UpdateTime: sgdUpdateTime,
		},
	}, nil
}
