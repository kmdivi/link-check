package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"
)

func ExtractLinks(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	re := regexp.MustCompile(`https://[^\s\)\]]+`)
	var links []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindAllString(line, -1)
		for _, match := range matches {
			links = append(links, match)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return links, nil
}

func CheckLink(url string) (int, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1788.0")

	resp, err := client.Do(req)

	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println("usage: ./link-check [markdown file]")
		os.Exit(1)
	}
	filePath := args[1]

	links, err := ExtractLinks(filePath)
	if err != nil {
		fmt.Println("Error extracting links:", err)
		return
	}

	for _, link := range links {
		statusCode, err := CheckLink(link)
		if err != nil {
			fmt.Printf("Error checking link %s: %v\n", link, err)
			continue
		}

		if statusCode == http.StatusOK {
			fmt.Printf("[Info] Link %s is accessible (status code: %d)\n", link, statusCode)
		} else {
			fmt.Printf("[Error] Link %s returned status code %d\n", link, statusCode)
		}
	}
}
