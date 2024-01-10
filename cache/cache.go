package cache

import (
	"github.com/go-redis/redis"
	"log"
)

type RedisCli struct {
	rdb *redis.Client
}

func InitRedis(url string, pwd string) (*RedisCli, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: pwd, // no password set
		DB:       0,   // use default DB
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &RedisCli{
		rdb: rdb,
	}, nil
}

func (r *RedisCli) ExistKey(key string) bool {
	n, err := r.rdb.Exists(key).Result()
	if err != nil {
		log.Println("find exist user Key error :", err)
	}
	if n == 0 {
		log.Println(key, "key no exist")
		return false
	}
	log.Println(key, "key exist")
	return true
}
