package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	discover := flag.Bool("discover", false, "enable discover mode")
	// Parse the command line flags
	flag.Parse()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		if input == "" {
			continue
		}

		// Resolve the IP address of the subdomain
		ip, err := net.LookupIP(input)
		if err != nil {
			fmt.Printf("Error resolving IP address for %s: %s\n", input, err)
			os.Exit(1)
		}

		// Check if the IP belongs to Cloudflare's range of IPs
		if isCloudflareIP(ip[0]) && *discover {
			fmt.Println(input)
		} else {
			if !isCloudflareIP(ip[0]) && !*discover {
				fmt.Println(input)
			}
		}

	}
}

func isIPInRange(ip net.IP, ipRanges []string) bool {
	for _, ipRange := range ipRanges {
		// Parse the IP range.
		_, network, err := net.ParseCIDR(ipRange)
		if err != nil {
			continue
		}

		// Check if the IP is in the range.
		if network.Contains(ip) {
			return true
		}
	}

	return false
}

// isCloudflareIP checks if the given IP belongs to Cloudflare's range of IPs
func isCloudflareIP(ip net.IP) bool {
	ipRanges := []string{
		"173.245.48.0/20",
		"103.21.244.0/22",
		"103.22.200.0/22",
		"103.31.4.0/22",
		"141.101.64.0/18",
		"108.162.192.0/18",
		"190.93.240.0/20",
		"188.114.96.0/20",
		"197.234.240.0/22",
		"198.41.128.0/17",
		"162.158.0.0/15",
		"104.16.0.0/13",
		"104.24.0.0/14",
		"172.64.0.0/13",
		"131.0.72.0/22",
	}
	if ip4 := ip.To4(); ip4 != nil {
		// Cloudflare's IPv4 range
		return isIPInRange(ip4, ipRanges)
	}

	// Cloudflare's IPv6 ranges are not yet covered
	return false
}
