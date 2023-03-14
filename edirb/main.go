package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	threads := flag.Int("t", 100, "Number of threads for requests")
	wordlist := flag.String("w", "", "Wordlist file path for bruteforce")

	flag.Parse()

	if *wordlist == "" {
		fmt.Fprintf(os.Stderr, "usage: %s -w wordlist target\n", os.Args[0])
		os.Exit(1)
	}

	// Read the items for bruteforcing from the wordlist file
	file, err := os.Open(*wordlist)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer file.Close()

	var items []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		item := scanner.Text()
		items = append(items, item)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading wordlist file:", err)
		os.Exit(1)
	}

	// Bruteforce items for each target concurrently
	var targets []string
	scanner = bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		target := scanner.Text()
		// Check if the target has a trailing slash
		if !strings.HasSuffix(target, "/") {
			target += "/"
		}
		targets = append(targets, target)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(1)
	}

	// Create a channel with buffer size 10 to limit threads
	sem := make(chan struct{}, *threads)

	for _, item := range items {
		var wg sync.WaitGroup
		for _, target := range targets {

			wg.Add(1)
			// Wait for a free slot in the channel
			sem <- struct{}{}

			go func(target, item string) {
				defer func() {
					// Release the slot in the channel
					<-sem
					wg.Done()
				}()

				url := fmt.Sprintf("%s%s", target, item)
				client := &http.Client{
					// so that you will not hang yourself over 120 seconds threads
					Timeout: time.Second * 5,
					CheckRedirect: func(req *http.Request, via []*http.Request) error {
						return http.ErrUseLastResponse
					},
				}

				resp, err := client.Get(url)
				if err != nil {
					return
				}
				var buf bytes.Buffer
				_, err = io.Copy(&buf, resp.Body)
				if err != nil {
					// handle error
					return
				}
				bodyLength := buf.Len()

				// this aproximates to the nearest hundreds so that you will not duplicate the outputs
				bodyLength = ((bodyLength + 50) / 100) * 100

				if resp.StatusCode != 404 {
					fmt.Printf("%s [sc:%d] [al:%d]\n", url, resp.StatusCode, bodyLength)
				}
				resp.Body.Close()
			}(target, item)
		}
		wg.Wait()
	}

}
