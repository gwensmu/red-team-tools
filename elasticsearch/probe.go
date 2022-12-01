package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func Login(host string) (ESCluster, error) {
	log.Println("Attempting to get cluster details for", host)

	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://"+host+":9200", nil)

	if err != nil {
		log.Fatal(err)
	}

	req.SetBasicAuth("elastic", "changeme")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	var es_cluster ESCluster
	es_cluster.Address = host

	error := json.NewDecoder(resp.Body).Decode(&es_cluster)

	if error != nil {
		log.Fatal(err)
	}

	resp.Body.Close()

	if resp.Status != "200 OK" {
		return es_cluster, errors.New("Login failed")
	}

	return es_cluster, nil
}
