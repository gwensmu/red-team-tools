package main

import (
	"fmt"
	"net"
)

const ES_DEFAULT_PORT = 9200
const NO_DICE = "No dice"

func worker(addresses []string, results chan string) {
	for ip := range addresses {
		address := fmt.Sprintf("%d:%d", ip, ES_DEFAULT_PORT)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- NO_DICE
			continue
		}
		conn.Close()
		results <- address
	}
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

	results := make(chan string)
	var public_instances []string

	for i := 0; i < cap(hosts); i++ {
		go worker(hosts, results)
	}

	for i := 0; i < cap(hosts); i++ {
		host := <-results
		if host != NO_DICE {
			public_instances = append(public_instances, host)
		}
	}

	close(results)
	for _, ip := range public_instances {
		fmt.Printf("%s running service on port %d\n", ip, ES_DEFAULT_PORT)
	}

	if len(public_instances) == 0 {
		fmt.Println("No public instances found")
	}
}
