package api

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"math/big"
	"metaoasis-filesystem/cache"
)

func (me *T) GetPriceInfoByAsset(hash string) ([]map[string]interface{}, error) {
	message := make(json.RawMessage, 0)
	ret := &message

	var currencyCode = [5]string{"KRW", "USD", "IDR", "SGD", "THB"}
	var r1, count, err = me.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "Asset",
		Index:      "getTotalsupply",
		Sort:       bson.M{},
		Filter:     bson.M{"hash": hash},
		Query:      []string{},
	}, ret)

	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0)
	totalsupply := big.NewInt(0)
	if count > 0 {
		decimal := r1[0]["decimals"].(int32)
		totalsupply, _, err = r1[0]["totalsupply"].(primitive.Decimal128).BigInt()
		totalsupply = totalsupply.Div(totalsupply, big.NewInt(int64(math.Pow10(int(decimal)))))
		if err != nil {
			return nil, err
		}
		// GET price from cache
		var gasPrice, neoPrice *cache.CurrencyList
		if r1[0]["symbol"].(string) == "GAS" {
			gasPrice, err = me.CacheClient.GetCacheGASPrice()
			if err != nil {
				return nil, err
			}
		} else if r1[0]["symbol"].(string) == "NEO" {
			neoPrice, err = me.CacheClient.GetCacheNeoPrice()
			if err != nil {
				return nil, err
			}
		}

		for i := 0; i < len(currencyCode); i++ {
			re := make(map[string]interface{})
			re["symbol"] = r1[0]["symbol"]
			re["currencyCode"] = currencyCode[i]
			re["price"] = nil
			re["marketCap"] = nil
			re["accTradePrice24h"] = nil
			re["circulatingSupply"] = totalsupply
			re["maxSupply"] = totalsupply
			re["provider"] = "NGD"
			re["lastUpdatedTimestamp"] = nil
			result = append(result, re)
			if re["symbol"].(string) == "GAS" {
				switch currencyCode[i] {
				case "KRW":
					re["price"] = gasPrice.KRW.Price
					re["lastUpdatedTimestamp"] = gasPrice.KRW.UpdateTime
				case "USD":
					re["price"] = gasPrice.USD.Price
					re["lastUpdatedTimestamp"] = gasPrice.USD.UpdateTime
				case "IDR":
					re["price"] = gasPrice.IDR.Price
					re["lastUpdatedTimestamp"] = gasPrice.IDR.UpdateTime
				case "SGD":
					re["price"] = gasPrice.SGD.Price
					re["lastUpdatedTimestamp"] = gasPrice.SGD.UpdateTime
				case "THB":
					re["price"] = gasPrice.THB.Price
					re["lastUpdatedTimestamp"] = gasPrice.THB.UpdateTime
				default:
					re["price"] = nil
				}
			} else if re["symbol"].(string) == "NEO" {
				switch currencyCode[i] {
				case "KRW":
					re["price"] = neoPrice.KRW.Price
					re["lastUpdatedTimestamp"] = neoPrice.KRW.UpdateTime
				case "USD":
					re["price"] = neoPrice.USD.Price
					re["lastUpdatedTimestamp"] = neoPrice.USD.UpdateTime
				case "IDR":
					re["price"] = neoPrice.IDR.Price
					re["lastUpdatedTimestamp"] = neoPrice.IDR.UpdateTime
				case "SGD":
					re["price"] = neoPrice.SGD.Price
					re["lastUpdatedTimestamp"] = neoPrice.SGD.UpdateTime
				case "THB":
					re["price"] = neoPrice.THB.Price
					re["lastUpdatedTimestamp"] = neoPrice.THB.UpdateTime
				default:
					re["price"] = nil
				}
			}

			if re["price"] != nil {
				re["marketCap"] = calculateMarketCap(totalsupply, re["price"])
			}

		}
	}

	return result, nil
}

func calculateMarketCap(maxSupplay *big.Int, price interface{}) interface{} {
	priceData := price.(float64)
	s := new(big.Float).SetInt(maxSupplay)
	muldata := s.Mul(s, new(big.Float).SetFloat64(priceData))

	return fmt.Sprintf("%.2f", muldata)
}
