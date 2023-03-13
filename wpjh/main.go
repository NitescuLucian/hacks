package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func logic(url_entered string) {

	//check for .com/ , remove the last / if it exists -> https://example.com/ -> https://example.com and we add /wp-json
	if strings.HasSuffix(url_entered, "/") {
		url_entered = strings.TrimSuffix(url_entered, "/")
	}

	url := strings.TrimSpace(url_entered)
	// Make HTTP GET request to the JSON API endpoint
	response, err := http.Get(url + "/wp-json")
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	// Parse JSON response
	var jsonData map[string]interface{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		panic(err)
	}

	// Extract links from JSON response
	links := extractLinks(jsonData)

	// Print links
	for _, link := range links {
		fmt.Println(link)
	}
}

func main() {

	/* If there is url supplied from cmdline, use it, if not then make use of stdin */
	if len(os.Args) > 1 {
		fmt.Println("Targeting -> ", os.Args[1], " recieved from command line")
		logic(os.Args[1])

	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			fmt.Println("Targeting -> ", scanner.Text(), " recieved from stdin")
			url := strings.TrimSpace(scanner.Text())
			logic(url)
		}
	}

}

func extractLinks(data map[string]interface{}) []string {
	links := []string{}

	// Recursively search for links in JSON object
	for _, value := range data {
		switch v := value.(type) {
		case []interface{}:
			for _, element := range v {
				if subData, ok := element.(map[string]interface{}); ok {
					subLinks := extractLinks(subData)
					links = append(links, subLinks...)
				}
			}
		case map[string]interface{}:
			subLinks := extractLinks(v)
			links = append(links, subLinks...)
		case string:
			// Check if string is a link
			if isLink(v) {
				links = append(links, v)
			}
		}
	}
	return links
}

// Helper function to check if a string is a valid URL
func isLink(s string) bool {
	_, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}
	u, err := url.Parse(s)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}
