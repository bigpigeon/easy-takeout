package cachemanage

import (
	"gopkg.in/redis.v5"
)

func CreateRedisClient(address string, password string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	return client, nil
}
