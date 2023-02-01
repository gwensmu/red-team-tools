package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net_helpers"
	"os"
	"probe"
)

const JUPYTER_DEFAULT_PORT = 8888

func GetAPIStatus(host string) (probe.PublicService, error) {
	log.Println("Attempting to get Jupyter Api status details for", host)

	var notebook = probe.PublicService{}
	notebook.Address = host

	req, err := http.NewRequest("GET", "http://"+host+":8888/api/status", nil)
	if err != nil {
		log.Print(err)
		return notebook, err
	}

	statusResponse, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	if statusResponse.StatusCode != 200 {
		return notebook, fmt.Errorf("jupyter status check failed with %d: %s", statusResponse.StatusCode, host)
	}

	body, error := ioutil.ReadAll(statusResponse.Body)
	if error != nil {
		fmt.Println(error)
	}

	defer statusResponse.Body.Close()

	log.Println(host + " has indexes:\n" + string(body))

	return notebook, nil
}

func Worker(addresses <-chan string, results chan probe.PublicService) {
	var nilInstance = probe.PublicService{}

	for ip := range addresses {

		_, err := net_helpers.Dial(ip, JUPYTER_DEFAULT_PORT)
		if err != nil {
			results <- nilInstance
			continue
		}

		instanceDetails, err := GetAPIStatus(ip)

		if err != nil {
			results <- nilInstance
			continue
		}

		log.Printf("Notebook %s is open\n", instanceDetails.Address)
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
