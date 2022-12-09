package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		url := strings.Split(strings.TrimSpace(scanner.Text()), "?")
		if len(url) > 1 {
			querry_params := strings.Split(url[1], "&")
			for _, item := range querry_params {
				fmt.Println(url[0] + "?" + item)
			}
		}
	}
}
