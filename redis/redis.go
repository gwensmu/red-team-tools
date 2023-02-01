package main

import (
	"context"
	"log"
	"net/http"
	"net_helpers"
	"os"
	"probe"

	"github.com/go-redis/redis/v8"
)

const REDIS_DEFAULT_PORT = 6379

type RedisInstance struct {
	Name    string
	Address string
}

var ctx = context.Background()

func GetKeys(host string) (probe.PublicService, error) {
	log.Println("Attempting to get instance details for", host)

	var redis_instance probe.PublicService
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
		log.Print(err)
	}

	return redis_instance, nil
}

func Worker(addresses <-chan string, results chan probe.PublicService) {
	var nilInstance = probe.PublicService{}

	for ip := range addresses {

		_, err := net_helpers.Dial(ip, REDIS_DEFAULT_PORT)
		if err != nil {
			results <- nilInstance
			continue
		}

		instanceDetails, err := GetKeys(ip)

		if err != nil {
			results <- nilInstance
			continue
		}

		log.Printf("Instance %s is open\n", instanceDetails.Address)
		results <- instanceDetails
	}
}

func main() {
	probe.InitLogFile("scans")

	cidrs_to_scan, workerCount, err := probe.ReadFlags()
	if err != nil {
		log.Fatal(err)
	}

	var workerHandler = probe.WorkerFunc(Worker)
	err = probe.Probe(cidrs_to_scan, workerCount, workerHandler)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
