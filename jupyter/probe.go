package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func GetAPIStatus(host string) (JupyterInstance, error) {
	log.Println("Attempting to get Jupyter Api status details for", host)

	var notebook = JupyterInstance{}
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
