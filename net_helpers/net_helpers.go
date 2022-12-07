package net_helpers

import (
	"fmt"
	"net"
	"time"
)

func Dial(ip string, port int) (string, error) {
	address := fmt.Sprintf("%s:%d", ip, port)

	d := net.Dialer{Timeout: 1 * time.Second}
	conn, err := d.Dial("tcp", address)

	if err != nil {
		return "", err
	}

	defer conn.Close()

	return "ok", nil
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

func main() {}
