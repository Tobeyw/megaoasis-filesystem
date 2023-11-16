package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"math/big"
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
		}
	}

	return result, nil
}