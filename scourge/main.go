package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/chromedp/chromedp"
)

func main() {
	// defer task.Release()
	// define a string flag
	var match string
	flag.StringVar(&match, "match", "", "a string match flag")

	// Create a new Chromedp context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// parse the flags
	flag.Parse()
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
