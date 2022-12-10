package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

func nucleitaglines(s string) string {
	// Compile a regular expression to match values between [ and ]
	re := regexp.MustCompile(`\[(.*?)\]`)

	// Retrieve all matches
	matches := re.FindAllStringSubmatch(s, -1)

	// Create a slice to hold the matched values
	values := make([]string, 0, len(matches))

	// Iterate over the matches and append the values to the slice
	for _, match := range matches {
		values = append(values, "["+match[1]+"]")
	}

	// Concatenate the values into a single string
	return strings.Join(values, " ")
}

// Function that detects if there are tags of the form "[ceva]" where "ceva" is random
func hasTags(s string) bool {
	// Compile a regular expression to match tags of the form "[ceva]"
	re := regexp.MustCompile(`\[(.*?)\]`)

	// Retrieve all matches
	matches := re.FindAllString(s, -1)

	// Return true if there are any matches, false otherwise
	return len(matches) > 0
}

// Function that extracts the domain and subdomain from a URL
func extractDomain(u string) (string, string) {
	// Parse the URL
	parsedURL, err := url.Parse(u)
	if err != nil {
		return "", ""
	}

	// Split the hostname into subdomains and the domain
	parts := strings.Split(parsedURL.Hostname(), ".")

	// Return the domain and the subdomain
	return parts[len(parts)-1], strings.Join(parts[:len(parts)-1], ".")
}

// Function that calculates the MD5 hash of a string and returns it as a string
func md5String(s string) string {
	// Calculate the MD5 hash of the string
	hash := md5.Sum([]byte(s))

	// Format the hash as a string and return it
	return fmt.Sprintf("%x", hash)
}

func main() {
	JIRAUSER := os.Getenv("JIRAUSER")
	JIRAAPI := os.Getenv("JIRAAPI")
	if JIRAUSER == "" || JIRAAPI == "" {
		fmt.Println("JIRAUSER or JIRAAPI is empty, provide valid creds")
		os.Exit(1)
	}
	var pid string
	var tag string
	var apiurl string
	scanner := bufio.NewScanner(os.Stdin)
	flag.StringVar(&pid, "pid", "", "The id of the jira board to which you want to report.")
	flag.StringVar(&tag, "tag", "", "The name tag you want to add to your finding.")
	flag.StringVar(&apiurl, "url", "", "The jira api url to which you are connected")
	flag.Parse()
	token := base64.StdEncoding.EncodeToString([]byte(JIRAUSER + ":" + JIRAAPI))
	burp0Headers := map[string]string{
		"Authorization": "Basic " + token,
		"User-Agent":    "curl/7.58.0",
		"Accept":        "*/*",
		"Content-Type":  "application/json",
		"Connection":    "close",
	}
	// Struct that represents the JSON body of the request
	type RequestBody struct {
		Fields struct {
			Description string `json:"description"`
			IssueType   struct {
				Name string `json:"name"`
			} `json:"issuetype"`
			Project struct {
				Key string `json:"key"`
			} `json:"project"`
			Summary string `json:"summary"`
		} `json:"fields"`
	}

	for scanner.Scan() {
		finding := scanner.Text()
		if finding == "" {
			continue
		}
		domain, subdomain := extractDomain(finding)
		titleinfo := "[" + subdomain + "." + domain + "]"
		if domain == "" && subdomain == "" {
			titleinfo = "[" + md5String(finding) + "]"
		}
		if hasTags(finding) {
			titleinfo = nucleitaglines(finding) + titleinfo
		}
		summary := "[" + tag + "] " + titleinfo
		if len(summary) > 200 {
			summary = "[" + md5String(finding) + "]"
		}
		burp0JSON := &RequestBody{
			Fields: struct {
				Description string `json:"description"`
				IssueType   struct {
					Name string `json:"name"`
				} `json:"issuetype"`
				Project struct {
					Key string `json:"key"`
				} `json:"project"`
				Summary string `json:"summary"`
			}{
				Description: "Check for " + tag + ":\n\n" + finding + "\n\nFinding in base64 escape:\n\n" + base64.StdEncoding.EncodeToString([]byte(finding)),
				IssueType: struct {
					Name string `json:"name"`
				}{
					Name: "Task",
				},
				Project: struct {
					Key string `json:"key"`
				}{
					Key: pid,
				},
				Summary: summary,
			},
		}
		// Marshal JSON body into a byte slice
		body, err := json.Marshal(burp0JSON)
		if err != nil {
			fmt.Println(err)
			return
		}
		// Create a new HTTP request
		req, err := http.NewRequest("POST", apiurl, bytes.NewBuffer(body))
		if err != nil {
			fmt.Println(err)
			return
		}
		// Set headers of the request
		for key, value := range burp0Headers {
			req.Header.Set(key, value)
		}
		check := false
		for !check {
			// Create a new HTTP client and do the request
			client := &http.Client{}
			res, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
				fmt.Println(finding)
				fmt.Println("Waiting...")
				time.Sleep(30 * time.Second)
				return
			}
			if res.StatusCode == 201 {
				check = true
			} else {
				fmt.Println(finding)
				fmt.Println("Waiting...")
				time.Sleep(30 * time.Second)
			}
		}
	}
}
