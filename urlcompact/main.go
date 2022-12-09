package main

import (
	"bufio"
	"crypto/md5"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	var fullscan bool
	flag.BoolVar(&fullscan, "full", false, "full mode ([^a-zA-Z0-9:-])")
	flag.Parse()
	lines := make(map[string]string)
	scanner := bufio.NewScanner(os.Stdin)
	// regex to check string for most commonly used letters of the English language and some numbers and ports ;)
	regex := regexp.MustCompile("^[etaionshru]|[etaionshru].*[etaionshru]$|[123895:-]$")
	if fullscan {
		regex = regexp.MustCompile("[^a-zA-Z0-9:]")
	}
	for scanner.Scan() {
		url := strings.Split(strings.TrimSpace(scanner.Text()), "?")
		querry_params := ""
		if len(url) > 1 {
			querry_params = url[1]
		}
		url_lookout := strings.Split(url[0], "/")
		if len(url_lookout) <= 3 {
			continue
		}
		paths_list := url_lookout[3:]
		paths_hash := md5.Sum([]byte(url_lookout[2] + regex.ReplaceAllString(strings.Join(paths_list, ""), "")))
		paths_hash_hex := fmt.Sprintf("%x", paths_hash)

		if len(lines[paths_hash_hex]) > 0 {
			if strings.Contains(lines[paths_hash_hex], "?") && len(querry_params) > 0 {
				selection_params := strings.Split(querry_params, "&")
				for _, param_value := range selection_params {
					key_value := strings.Split(param_value, "=")[0]
					if strings.Index(lines[paths_hash_hex], key_value+"=") == -1 {
						lines[paths_hash_hex] = lines[paths_hash_hex] + "&" + param_value
					}
				}
			}
			continue
		}

		if len(querry_params) == 0 {
			lines[paths_hash_hex] = url[0]
			continue
		}
		lines[paths_hash_hex] = url[0] + "?" + querry_params
	}
	for _, str := range lines {
		fmt.Println(str)
	}
}
