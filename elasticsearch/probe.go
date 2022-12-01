package main

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aquasecurity/esquery"
	"github.com/elastic/go-elasticsearch/v7"
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

	// todo: better default query
	GetIndexes(es)

	return es_cluster, nil
}

func GetIndexes(es *elasticsearch.Client) {
	// run a boolean search query
	res, err := esquery.Search().
		Query(
			esquery.
				Bool().
				Must(esquery.Term("title", "Go and Stuff")).
				Filter(esquery.Term("tag", "tech")),
		).
		Aggs(
			esquery.Avg("average_score", "score"),
			esquery.Max("max_score", "score"),
		).
		Size(20).
		Run(
			es,
			es.Search.WithContext(context.TODO()),
			es.Search.WithIndex("test"),
		)
	if err != nil {
		log.Fatalf("Failed searching for stuff: %s", err)
	}

	// ensure that we conseme the response body
	io.Copy(ioutil.Discard, res.Body)
	defer res.Body.Close()
}
