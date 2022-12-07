package net_helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

type AWSPrefix struct {
	IPPrefix string `json:"ip_prefix"`
	Service  string `json:"service"`
	Region   string `json:"region"`
}

type AWSPrefixes struct {
	Prefixes []AWSPrefix `json:"prefixes"`
}

func GetAWSCIDRs(region string) []string {
	req, err := http.NewRequest("GET", "https://ip-ranges.amazonaws.com/ip-ranges.json", nil)
	if err != nil {
		log.Fatal(err)
	}

	res, _ := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	var blocks AWSPrefixes
	d := json.NewDecoder(res.Body)
	d.Decode(&blocks)

	if err != nil {
		log.Fatal(err)
	}

	var prefixes []string

	for _, prefix := range blocks.Prefixes {
		if prefix.Region == region {
			prefixes = append(prefixes, prefix.IPPrefix)
		}
	}

	return prefixes
}

type GCEPrefix struct {
	IPv4Prefix string `json:"ipv4Prefix"`
	IPv6Prefix string `json:"ipv6Prefix"`
	Service    string `json:"service"`
	Scope      string `json:"scope"`
}

type GCEPrefixes struct {
	Prefixes []GCEPrefix `json:"prefixes"`
}

func GetGCEPrefixes(region string) []string {
	req, err := http.NewRequest("GET", "https://www.gstatic.com/ipranges/cloud.json", nil)

	if err != nil {
		log.Fatal(err)
	}

	res, _ := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	var blocks GCEPrefixes
	d := json.NewDecoder(res.Body)
	d.Decode(&blocks)

	if err != nil {
		log.Fatal(err)
	}

	var prefixes []string

	// todo: filter on region first for clarity
	for _, prefix := range blocks.Prefixes {
		if prefix.Scope == region && prefix.IPv4Prefix != "" {
			prefixes = append(prefixes, prefix.IPv4Prefix)
		} else if prefix.Scope == region && prefix.IPv6Prefix != "" {
			prefixes = append(prefixes, prefix.IPv6Prefix)
		}
	}

	return prefixes
}

func GetCIDR(cloud string, region string) []string {
	var cidrs []string

	switch cloud {
	case "aws":
		cidrs = GetAWSCIDRs(region)
	case "gce":
		cidrs = GetGCEPrefixes(region)
	default:
		log.Fatalf("Cloud provider %s not supported", cloud)
	}

	return cidrs
}
