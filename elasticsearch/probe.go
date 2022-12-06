package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

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
			"https://" + host + ":9200",
			"https://" + host + ":9201",
		},
		Username: "elastic",
		Password: "changeme",
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed creating client: %s", err)
	}

	GetIndexes(es)

	return es_cluster, nil
}

func GetIndexes(es *elasticsearch.Client) {
	res, err := esapi.CatIndicesRequest{Format: "json"}.Do(context.Background(), es)
	if err != nil {
		return
	}
	defer res.Body.Close()

	fmt.Println(res.String())
}
