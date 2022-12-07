package main

import (
	"log"
	"net/http"
)

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

	// todo: print keys of redis instance

	return redis_instance, nil
}
