package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		if input == "" {
			continue
		}
		output, err := url.QueryUnescape(input)
		if err != nil {
			continue
		}
		fmt.Println(output)
	}
}
