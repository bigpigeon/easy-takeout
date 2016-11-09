package cachemanage

import (
	"gopkg.in/redis.v5"
)

type Client struct {
	redis.Client
}

func CreateClient(address string) *Client {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       1,
	})
	return &Client{*client}
}
