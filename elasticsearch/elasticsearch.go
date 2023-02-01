package main

import (
	"context"
	"fmt"
	"log"
	"net_helpers"
	"os"
	"probe"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

const ES_DEFAULT_PORT = 9200

// can you imagine? years of this
const ES_DEFAULT_USERNAME = "elastic"
const ES_DEFAULT_PASSWORD = "changeme"

func Login(host string) (probe.PublicService, error) {
	log.Println("Attempting to get cluster details for", host)

	var es_cluster probe.PublicService
	es_cluster.Address = host

	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://" + host + ":9200",
			"http://" + host + ":9201",
		},
		Username: ES_DEFAULT_USERNAME,
		Password: ES_DEFAULT_PASSWORD,
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed creating client: %s", err)
		return es_cluster, err
	}

	indexes, err := GetIndexes(es, host)
	if err != nil {
		log.Printf("Failed getting indexes: %s", err)
		return es_cluster, err
	} else {
		log.Println(host + " has indexes:\n" + indexes)
		return es_cluster, nil
	}
}

func GetIndexes(es *elasticsearch.Client, host string) (string, error) {
	res, err := esapi.CatIndicesRequest{Pretty: true}.Do(context.Background(), es)
	if err != nil {
		return fmt.Sprintf("Error getting indexes: %s", err), err
	}

	if res.Status() == "401" {
		return fmt.Sprintf("401 Unauthorized: %s", host), err
	}

	defer res.Body.Close()

	return res.String(), err
}

func Worker(addresses <-chan string, results chan probe.PublicService) {
	var nilCluster = probe.PublicService{}

	for ip := range addresses {

		_, err := net_helpers.Dial(ip, ES_DEFAULT_PORT)
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
