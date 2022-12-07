package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func GetKeys(host string) (RedisInstance, error) {
	log.Println("Attempting to get instance details for", host)

	var redis_instance RedisInstance
	redis_instance.Address = host

	// todo: magic number with port
	_, err := http.NewRequest("GET", "http://"+host+":6379", nil)
	if err != nil {
		log.Print(err)
		return redis_instance, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	iter := rdb.Scan(ctx, 0, "prefix:*", 0).Iterator()
	for iter.Next(ctx) {
		log.Println("keys", iter.Val())
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}

	return redis_instance, nil
}
