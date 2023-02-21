package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/big"
)

func (me *T) GetTotalSupplyByAsset(hash string) (*big.Int, error) {
	message := make(json.RawMessage, 0)
	ret := &message

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

	totalsupply := big.NewInt(0)
	if count > 0 {
		totalsupply, _, err = r1[0]["totalsupply"].(primitive.Decimal128).BigInt()
		if err != nil {
			return nil, err
		}
	}

	return totalsupply, nil
}
