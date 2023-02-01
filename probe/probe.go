package probe

import (
	"flag"
	"fmt"
	"log"
	"net_helpers"
	"os"
	"time"
)

type probe interface {
	InitLogFile() error
	ReadFlags() ([]string, int, error)
	Probe() error
	Worker()
}
type PublicService struct {
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

type WorkerFunc func(addresses <-chan string, results chan PublicService)

func Probe(cidrs []string, workerCount int, w WorkerFunc) error {
	for _, block := range cidrs {
		hosts, _ := net_helpers.Hosts(block)

		log.Println("Scanning", len(hosts), "hosts in CIDR", block)

		addresses := make(chan string, len(hosts))
		for _, host := range hosts {
			addresses <- host
		}

		results := make(chan PublicService)
		var public_instances []PublicService

		for i := 0; i < workerCount; i++ {
			go w(addresses, results)
		}

		close(addresses)

		for i := 0; i < len(hosts); i++ {
			instance := <-results

			if instance.Name != "" {
				public_instances = append(public_instances, instance)
			}
		}

		close(results)

		fmt.Println("Found", len(public_instances), "public services")
	}

	return nil
}

func ReadFlags() ([]string, int, error) {
	blockPtr := flag.String("block", "", "a IPv4 CIDR block to scan")
	cloudProviderPtr := flag.String("cloud", "aws", "the cloud provider to scan (aws/gce)")
	regionPtr := flag.String("region", "us-east1", "the region to scan")
	workerPtr := flag.Int("workers", 20, "the number of workers to use")

	flag.Parse()

	var cidrs_to_scan []string

	if *blockPtr != "" {
		cidrs_to_scan = []string{*blockPtr}
	} else {
		cidrs_to_scan = net_helpers.GetCIDR(*cloudProviderPtr, *regionPtr)
	}

	return cidrs_to_scan, *workerPtr, nil
}

func InitLogFile(dir string) error {
	filename := fmt.Sprintf("%s/scan-%s.log", dir, time.Now())
	logFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(logFile)

	return err
}
