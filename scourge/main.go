package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	// defer task.Release()
	// define a string flag
	var match string
	var maxtime int
	flag.StringVar(&match, "match", "", "a string match flag")
	flag.IntVar(&maxtime, "maxtime", 120, "a string to set seconds of maximum execution")
	// parse the flags
	flag.Parse()

	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Duration(maxtime)*time.Second)
	defer cancel()

	// Create a new Chromedp context with the timeout context
	ctx, cancel := chromedp.NewContext(timeoutCtx)
	defer cancel()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		if input == "" {
			continue
		}
		// Navigate to a page with query parameters
		err := chromedp.Run(ctx,
			chromedp.Navigate(input),
		)
		if err != nil {
			continue
		}

		var body string
		err = chromedp.Run(ctx, chromedp.OuterHTML("html", &body))
		if err != nil {
			continue
		}
		if strings.Contains(body, match) && len(match) > 0 {
			fmt.Println(input)
		}

	}
}
