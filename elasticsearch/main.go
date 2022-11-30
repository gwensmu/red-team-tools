package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

const ES_DEFAULT_PORT = 9200
const NO_DICE = "No dice"

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

		_, err := Dial(ip)
		if err != nil {
			results <- nilCluster
			continue
		}

		clusterDetails, err := Login(ip)

		if err != nil {
			results <- nilCluster
			continue
		}

		results <- clusterDetails
	}
}

func Dial(ip string) (string, error) {
	log.Println("Attempting to dial", ip)
	address := fmt.Sprintf("%s:%d", ip, ES_DEFAULT_PORT)

	d := net.Dialer{Timeout: 1 * time.Second}
	conn, err := d.Dial("tcp", address)

	if err != nil {
		return "", err
	}

	defer conn.Close()

	return "ok", nil
}

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

func Hosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	// remove network address and broadcast address
	return ips[1 : len(ips)-1], nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func main() {
	var cidr_to_scan string

	fmt.Print("Enter CIDR to scan: ")
	fmt.Scanln(&cidr_to_scan)

	hosts, _ := Hosts(cidr_to_scan)

	fmt.Println("Scanning", len(hosts), "hosts")

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
