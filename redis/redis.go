package main

import (
	"flag"
	"fmt"
	"log"
	"net_helpers"
	"os"
	"time"
)

const REDIS_DEFAULT_PORT = 6379
const NO_DICE = "No dice"

var logFileDir = "scans"

type RedisInstance struct {
	Name    string
	Address string
}

func worker(addresses <-chan string, results chan RedisInstance) {
	var nilInstance = RedisInstance{}

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

func initLogFile(dir string) {
	filename := fmt.Sprintf("%s/redis-scan-%s.log", dir, time.Now())
	logFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(logFile)
}

func main() {
	initLogFile(logFileDir)

	blockPtr := flag.String("block", "", "a IPv4 CIDR block to scan")
	cloudProviderPtr := flag.String("cloud", "aws", "the cloud provider to scan (aws/gce)")
	regionPtr := flag.String("region", "us-east1", "the region to scan")

	flag.Parse()

	var cidrs_to_scan []string

	if *blockPtr != "" {
		cidrs_to_scan = []string{*blockPtr}
	} else {
		cidrs_to_scan = net_helpers.GetCIDR(*cloudProviderPtr, *regionPtr)
	}

	for _, block := range cidrs_to_scan {
		hosts, _ := net_helpers.Hosts(block)

		log.Println("Scanning", len(hosts), "hosts in CIDR", block)

		addresses := make(chan string, len(hosts))
		for _, host := range hosts {
			addresses <- host
		}

		results := make(chan RedisInstance)
		var public_instances []RedisInstance

		for i := 0; i < 20; i++ {
			go worker(addresses, results)
		}

		close(addresses)

		for i := 0; i < len(hosts); i++ {
			instance := <-results

			if instance.Name != "" {
				public_instances = append(public_instances, instance)
			}
		}

		close(results)

		fmt.Println("Found", len(public_instances), "public redis instances")
	}

	os.Exit(0)
}
