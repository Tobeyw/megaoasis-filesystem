package api

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"metaoasis-filesystem/consts"
	"os"
)

func (me *T) GetMarketWhiteList() ([]string, error) {
	message := make(json.RawMessage, 0)
	ret := &message

	rt := os.ExpandEnv("${RUNTIME}")
	market := consts.Market_Main
	switch rt {
	case "test":
		market = consts.Market_Test
	case "staging":
		market = consts.Market_Main
	default:
		fmt.Sprintf("runtime environment mismatch")
	}
	pipeline := []bson.M{
		bson.M{"$match": bson.M{"market": market}},
		bson.M{"$match": bson.M{"$or": []interface{}{
			bson.M{"eventname": "AddAsset"},
			bson.M{"eventname": "RemoveAsset"},
		}}},
		bson.M{"$sort": bson.M{"timestamp": 1}},
		bson.M{"$group": bson.M{"_id": "$asset", "asset": bson.M{"$last": "$asset"}, "eventname": bson.M{"$last": "$eventname"}}},
	}

	var r1, err = me.Client.QueryAggregate(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Pipeline   []bson.M
			Query      []string
		}{
			Collection: "MarketNotification",
			Index:      "GetNFTClass",
			Sort:       bson.M{},
			Filter:     bson.M{},
			Pipeline:   pipeline,
			Query:      []string{},
		}, ret)

	if err != nil {
		return nil, err
	}

	var assetArr []string
	for _, item := range r1 {
		if item["eventname"].(string) == "AddAsset" {
			assetArr = append(assetArr, item["asset"].(string))
		}
	}
	return assetArr, nil
}
