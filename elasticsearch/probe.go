package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

const ES_DEFAULT_USERNAME = "elastic"
const ES_DEFAULT_PASSWORD = "changeme"

func Login(host string) (ESCluster, error) {
	log.Println("Attempting to get cluster details for", host)

	var es_cluster ESCluster
	es_cluster.Address = host

	_, err := http.NewRequest("GET", "http://"+host+":9200", nil)
	if err != nil {
		log.Print(err)
		return es_cluster, err
	}

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
