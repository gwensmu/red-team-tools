package probe

type probe interface {
	InitLogFile() error
	ReadFlags() []string, int, error
	Probe() error
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

func (cidrs []string, workerCount int, ) Probe() error {
	for _, block := range cidrs_to_scan {
		hosts, _ := net_helpers.Hosts(block)

		log.Println("Scanning", len(hosts), "hosts in CIDR", block)

		addresses := make(chan string, len(hosts))
		for _, host := range hosts {
			addresses <- host
		}

		results := make(chan PublicService)
		var public_instances []PublicService

		for i := 0; i < *workerPtr; i++ {
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

		fmt.Println("Found", len(public_instances), "public services")
	}
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
