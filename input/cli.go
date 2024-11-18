package input

import (
	"flag"
	"fmt"
)

func Parse() (string, error) {
	var url string
	flag.StringVar(&url, "url", "", "youtube video url")
	flag.Parse()

	if url == "" {
		return "", fmt.Errorf("url is empty")
	}

	return url, nil
}
