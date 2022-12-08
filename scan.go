package main

import (
	"flag"
	"fmt"
	"log"
	"net_helpers"
	"os"
	es "red-team-tools/elasticsearch"
	jupyter "red-team-tools/jupyter"
	redis "red-team-tools/redis"
	"time"
)

const JUPYTER_DEFAULT_PORT = 8888
const ES_DEFAULT_PORT = 9200
const REDIS_DEFAULT_PORT = 6379
const NO_DICE = "No dice"

var logFileDir = "scans"

var esLog = log.Default()
var redisLog = log.Default()
var jupyterLog = log.Default()
var mainLog = log.Default()

func esScan(ip string, esResults chan es.ESCluster) {
	var nilESCluster = es.ESCluster{}

	_, err := net_helpers.Dial(ip, ES_DEFAULT_PORT)
	if err != nil {
		esResults <- nilESCluster
	} else {
		clusterDetails, err := es.Login(ip)

		if err != nil {
			esResults <- nilESCluster
		}

		esLog.Printf("cluster %s (v%s) is open (%s)\n", clusterDetails.Cluster_Name, clusterDetails.Version.Number, clusterDetails.Address)
		esResults <- clusterDetails
	}
}

func redisScan(ip string, redisResults chan redis.RedisInstance) {
	var nilRedisInstance = redis.RedisInstance{}

	_, err := net_helpers.Dial(ip, REDIS_DEFAULT_PORT)
	if err != nil {
		redisResults <- nilRedisInstance
	} else {
		instanceDetails, err := redis.GetKeys(ip)

		if err != nil {
			redisResults <- nilRedisInstance
		}

		redisLog.Printf("Instance %s is open\n", instanceDetails.Address)
		redisResults <- instanceDetails
	}
}

func jupyterScan(ip string, jupyterResults chan jupyter.JupyterInstance) {
	var nilJupyterInstance = jupyter.JupyterInstance{}

	_, err := net_helpers.Dial(ip, JUPYTER_DEFAULT_PORT)
	if err != nil {
		jupyterResults <- nilJupyterInstance
	} else {
		instanceDetails, err := jupyter.GetAPIStatus(ip)

		if err != nil {
			jupyterResults <- nilJupyterInstance
		}

		jupyterLog.Printf("Notebook %s is open\n", instanceDetails.Address)
		jupyterResults <- instanceDetails
	}
}

func worker(addresses <-chan string, esResults chan es.ESCluster, redisResults chan redis.RedisInstance, jupyterResults chan jupyter.JupyterInstance) {
	for ip := range addresses {
		esScan(ip, esResults)
		redisScan(ip, redisResults)
		jupyterScan(ip, jupyterResults)
	}
}

func initLogFile(dir string, service string) *os.File {
	timestamp := time.Now().Format("01-01-2006-15-04-05")
	os.Mkdir("scans/"+timestamp+"", 0777)
	filename := fmt.Sprintf("%s/%s/%s-scan.log", dir, timestamp, service)
	logFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		mainLog.Fatalf("error opening file: %v", err)
	}

	return logFile
}

func main() {
	esLog.SetOutput(initLogFile(logFileDir, "elasticsearch"))
	redisLog.SetOutput(initLogFile(logFileDir, "redis"))
	jupyterLog.SetOutput(initLogFile(logFileDir, "jupyter"))
	mainLog.SetOutput(initLogFile(logFileDir, "main"))

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

		mainLog.Println("Scanning", len(hosts), "hosts in CIDR", block)

		addresses := make(chan string, len(hosts))
		for _, host := range hosts {
			addresses <- host
		}

		esResults := make(chan es.ESCluster)
		redisResults := make(chan redis.RedisInstance)
		jupyterResults := make(chan jupyter.JupyterInstance)

		for i := 0; i < 20; i++ {
			go worker(addresses, esResults, redisResults, jupyterResults)
		}

		close(addresses)
		close(esResults)
		close(jupyterResults)
		close(redisResults)
	}

	os.Exit(0)
}
