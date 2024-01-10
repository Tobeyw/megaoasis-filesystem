package cache

import (
	"encoding/json"
	"log"
)

const (
	CACHE_KEY_NEO = "CACHE_NEO"
	CACHE_KEY_GAS = "CACHE_GAS"
)

type Price struct {
	KRW interface{}
	USD interface{}
	IDR interface{}
	SGD interface{}
	THB interface{}
}

func (r *RedisCli) SetCacheNeoPrice(data Price) (err error) {
	strdata, _ := json.Marshal(data)
	err = r.rdb.Set(CACHE_KEY_NEO, strdata, 0).Err()
	if err != nil {
		log.Println("SET redis CACHE_NEO error", err)
		return err
	}
	return nil
}

func (r *RedisCli) SetCacheGASPrice(data Price) (err error) {
	strdata, _ := json.Marshal(data)
	err = r.rdb.Set(CACHE_KEY_GAS, strdata, 0).Err()
	if err != nil {
		log.Println("SET redis CACHE_GAS error", err)
		return err
	}
	return nil
}

func (r *RedisCli) GetCacheNeoPrice() (data *Price, err error) {
	res, err := r.rdb.Get(CACHE_KEY_NEO).Result()
	if err != nil {
		log.Println("GET redis neo  price error:", err)
		return nil, err
	}
	err = json.Unmarshal([]byte(res), &data)
	return data, nil
}

func (r *RedisCli) GetCacheGASPrice() (data *Price, err error) {
	res, err := r.rdb.Get(CACHE_KEY_GAS).Result()
	if err != nil {
		log.Println("GET redis gas price error:", err)
		return nil, err
	}
	err = json.Unmarshal([]byte(res), &data)
	return data, nil
}
