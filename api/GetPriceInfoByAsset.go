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
		var gasPrice, neoPrice *cache.Price
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
					re["price"] = gasPrice.KRW
				case "USD":
					re["price"] = gasPrice.USD
				case "IDR":
					re["price"] = gasPrice.IDR
				case "SGD":
					re["price"] = gasPrice.SGD
				case "THB":
					re["price"] = gasPrice.THB
				default:
					re["price"] = nil
				}
			} else if re["symbol"].(string) == "NEO" {
				switch currencyCode[i] {
				case "KRW":
					re["price"] = neoPrice.KRW
				case "USD":
					re["price"] = neoPrice.USD
				case "IDR":
					re["price"] = neoPrice.IDR
				case "SGD":
					re["price"] = neoPrice.SGD
				case "THB":
					re["price"] = neoPrice.THB
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
