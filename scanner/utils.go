package scanner

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadFile(path string) ([]string, error) {
	fileHandler, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	defer fileHandler.Close()

	var urls []string
	scanner := bufio.NewScanner(fileHandler)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url != "" {
			urls = append(urls, url)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return urls, nil
}
