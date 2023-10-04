package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func main() {
	lines := make(map[string]string)
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		input := scanner.Text()
		if input == "" {
			continue
		}
		if _, ok := lines[input]; !ok {
			lines[input] = input
			extractParametersFromURL(input)
		}
	}
}

func extractParametersFromURL(url string) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Ignore certificate errors
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(allocCtx)
	defer cancel()

	var htmlContent string

	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.OuterHTML("html", &htmlContent),
	); err != nil {
		return
	}

	parameters := extractParametersFromString(htmlContent)
	for name, value := range parameters {
		fmt.Printf("%s?%s=%s\n", url, name, value)
	}
}

func extractParametersFromString(input string) map[string]string {
	parameters := make(map[string]string)

	// Define regular expressions
	reName := regexp.MustCompile(`(?i)name=("|')?`)
	reInputs := regexp.MustCompile(`(?i)name=("|')?[\w-]+"`)
	reVar := regexp.MustCompile(`(?i)(var|let|const)\s+?`)
	reFullVars := regexp.MustCompile(`(?i)(var|let|const)\s+?[\w-]+"`)
	reWordsInQuotes := regexp.MustCompile(`("|')[a-zA-Z0-9]{3,20}('|")`)
	reWordsWithinObjects := regexp.MustCompile(`[\{,]\s*[[:alpha:]]\w{2,25}:#`)

	// Process the input string with regular expressions
	for _, re := range []*regexp.Regexp{reName, reInputs, reFullVars, reWordsInQuotes, reWordsWithinObjects, reVar} {
		matches := re.FindAllString(input, -1)
		for _, match := range matches {
			// Strip quotes (" or '), spaces, and equal sign (=) from the matched substring
			cleanedMatch := strings.Map(func(r rune) rune {
				if r == '"' || r == '\'' || r == '=' || unicode.IsSpace(r) {
					return -1 // Remove quotes, spaces, and equal sign
				}
				return r
			}, match)

			parameters[cleanedMatch] = "ph"
		}
	}

	// Create a goquery document from the input string
	reader := strings.NewReader(input)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return parameters
	}

	// Define a list of HTML element types to search for with the "name" attribute
	elementTypes := []string{"input", "button", "iframe", "fieldset", "meta", "output", "param", "select", "textarea"}

	// Loop through each element type
	for _, elementType := range elementTypes {
		doc.Find(elementType + "[name]").Each(func(index int, item *goquery.Selection) {
			nameAttr, exists := item.Attr("name")
			if exists {
				valueAttr, _ := item.Attr("value")
				// Extract value enclosed in single or double quotes
				value := extractQuotedValue(valueAttr)
				parameters[nameAttr] = value
			}
		})
	}

	// Find JavaScript variables in script tags and extract them - not performing well -
	doc.Find("script").Each(func(index int, script *goquery.Selection) {
		text := script.Text()
		lines := strings.Split(text, "\n")
		for _, line := range lines {
			if strings.Contains(line, "var ") {
				parts := strings.Split(line, "var ")
				if len(parts) == 2 {
					nameValue := strings.Split(parts[1], "=")
					if len(nameValue) == 2 {
						name := strings.TrimSpace(nameValue[0])
						parameters[name] = "ph"
					}
				}
			}
		}
	})

	return parameters
}

// Extracts the value enclosed in single or double quotes
func extractQuotedValue(input string) string {
	if strings.HasPrefix(input, "'") && strings.HasSuffix(input, "'") {
		return strings.Trim(input, "'")
	} else if strings.HasPrefix(input, "\"") && strings.HasSuffix(input, "\"") {
		return strings.Trim(input, "\"")
	}
	return input
}
