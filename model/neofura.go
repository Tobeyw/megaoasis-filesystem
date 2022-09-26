package model

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type T struct {
	Db_online string
	C_online  *mongo.Client
	Ctx       context.Context
}

func (me *T) QueryAll(args struct {
	Collection string
	Index      string
	Sort       bson.M
	Filter     bson.M
	Query      []string
	Limit      int64
	Skip       int64
}, ret *json.RawMessage) ([]map[string]interface{}, int64, error) {
	var results []map[string]interface{}
	convert := make([]map[string]interface{}, 0)
	collection := me.C_online.Database(me.Db_online).Collection(args.Collection)
	op := options.Find()
	op.SetSort(args.Sort)
	op.SetLimit(args.Limit)
	op.SetSkip(args.Skip)
	co := options.CountOptions{}
	count, err := collection.CountDocuments(me.Ctx, args.Filter, &co)
	if err != nil {
		return nil, 0, err
	}
	cursor, err := collection.Find(me.Ctx, args.Filter, op)
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			//log2.Fatalf("Closing cursor error %v", err)
			fmt.Println("Closing cursor error %v", err)
		}
	}(cursor, me.Ctx)
	if err == mongo.ErrNoDocuments {
		return nil, 0, err
	}
	if err != nil {
		return nil, 0, err
	}
	if err = cursor.All(me.Ctx, &results); err != nil {
		return nil, 0, err
	}
	for _, item := range results {
		if len(args.Query) == 0 {
			convert = append(convert, item)
		} else {
			temp := make(map[string]interface{})
			for _, v := range args.Query {
				temp[v] = item[v]
			}
			convert = append(convert, temp)
		}
	}
	r, err := json.Marshal(convert)
	if err != nil {
		return nil, 0, err
	}
	*ret = json.RawMessage(r)
	return convert, count, nil
}


func (me *T) QueryAggregate(args struct {
	Collection string
	Index      string
	Sort       bson.M
	Filter     bson.M
	Pipeline   []bson.M
	Query      []string
}, ret *json.RawMessage) ([]map[string]interface{}, error) {

	var results []map[string]interface{}
	convert := make([]map[string]interface{}, 0)
	collection := me.C_online.Database(me.Db_online).Collection(args.Collection)
	op := options.AggregateOptions{}

	cursor, err := collection.Aggregate(me.Ctx, args.Pipeline, &op)

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Fatalf("Closing cursor error %v", err)
		}
	}(cursor, me.Ctx)
	if err == mongo.ErrNoDocuments {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	if err = cursor.All(me.Ctx, &results); err != nil {
		return nil, err
	}

	for _, item := range results {
		if len(args.Query) == 0 {
			convert = append(convert, item)
		} else {
			temp := make(map[string]interface{})
			for _, v := range args.Query {
				temp[v] = item[v]
			}
			convert = append(convert, temp)
		}
	}

	r, err := json.Marshal(convert)
	if err != nil {
		return nil, err
	}
	*ret = json.RawMessage(r)
	return convert, nil
}