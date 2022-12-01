package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

const ES_DEFAULT_PORT = 9200
const NO_DICE = "No dice"

var logFileDir = "scans"

type ESCluster struct {
	Name         string
	Address      string
	Cluster_Name string
	ClusterUuid  string
	Version      struct {
		Number      string
		BuildFlavor string
		BuildType   string
	}
}

func worker(addresses <-chan string, results chan ESCluster) {
	var nilCluster = ESCluster{}

	for ip := range addresses {

		_, err := Dial(ip, ES_DEFAULT_PORT)
		if err != nil {
			results <- nilCluster
			continue
		}

		clusterDetails, err := Login(ip)

		if err != nil {
			results <- nilCluster
			continue
		}

		log.Printf("cluster %s (v%s) is open (%s)\n", clusterDetails.Cluster_Name, clusterDetails.Version.Number, clusterDetails.Address)
		results <- clusterDetails
	}
}

func initLogFile(dir string) {
	filename := fmt.Sprintf("%s/elasticsearch-scan-%s.log", dir, time.Now())
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
		cidrs_to_scan = GetCIDR(*cloudProviderPtr, *regionPtr)
	}

	for _, block := range cidrs_to_scan {
		hosts, _ := Hosts(block)

		log.Println("Scanning", len(hosts), "hosts in CIDR", block)

		addresses := make(chan string, len(hosts))
		for _, host := range hosts {
			addresses <- host
		}

		results := make(chan ESCluster)
		var public_instances []ESCluster

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

		fmt.Println("Found", len(public_instances), "public instances")

		for _, instance := range public_instances {
			log.Printf("cluster %s (v%s) is open (%s)\n", instance.Cluster_Name, instance.Version.Number, instance.Address)
		}

		if len(public_instances) == 0 {
			log.Println("No public instances found")
		}
	}

	os.Exit(0)
}
