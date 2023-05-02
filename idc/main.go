package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		input := scanner.Text()
		if input == "" {
			continue
		}

		// Get the IP address of the domain name
		ip, err := net.LookupIP(input)
		if err != nil {
			continue
		}
		// Query Shodan for the IP address
		resp, err := http.Get(fmt.Sprintf("https://internetdb.shodan.io/%s", ip[0].String()))
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		// Parse the JSON response
		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			continue
		}
		// Print the vuln if they exist
		if vulns, ok := result["vulns"]; ok {
			for _, vuln := range vulns.([]interface{}) {
				fmt.Printf("%s:%s\n", input, vuln)
			}
		}
	}

}
