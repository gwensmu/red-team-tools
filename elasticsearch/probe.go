package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

func Login(host string) (ESCluster, error) {
	log.Println("Attempting to get cluster details for", host)

	var es_cluster ESCluster
	es_cluster.Address = host

	client := &http.Client{
		Timeout: time.Second * 1,
	}

	req, err := http.NewRequest("GET", "http://"+host+":9200", nil)
	if err != nil {
		log.Print(err)
		return es_cluster, err
	}

	req.SetBasicAuth("elastic", "changeme")
	response, err := client.Do(req)
	if err != nil {
		log.Printf("Basic Auth failed for host %s: %s", host, err)
		return es_cluster, err
	}
	defer response.Body.Close()

	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://" + host + ":9200",
			"http://" + host + ":9201",
		},
		Username: "elastic",
		Password: "changeme",
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed creating client: %s", err)
	}

	indexes := GetIndexes(es)
	log.Println(host + " has indexes:\n" + indexes)

	return es_cluster, nil
}

func GetIndexes(es *elasticsearch.Client) string {
	res, err := esapi.CatIndicesRequest{Pretty: true}.Do(context.Background(), es)
	if err != nil {
		return fmt.Sprintf("Error getting indexes: %s", err)
	}
	defer res.Body.Close()

	return res.String()
}
