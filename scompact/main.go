package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {

	lines := make(map[string]string)
	scanner := bufio.NewScanner(os.Stdin)
	regex := regexp.MustCompile("[etaionshru]+")

	for scanner.Scan() {
		input := scanner.Text()
		if input == "" {
			continue
		}
		str := strings.Join(regex.FindAllString(input, -1), "")
		if _, ok := lines[str]; !ok {
			lines[str] = input
			fmt.Println(input)
		}

	}
}
