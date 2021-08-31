package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	for _, url := range getTrending() {
		fmt.Printf("%s | %s\n", url, getAbout(url))
	}
}

func getTrending() []string {
	var results []string
	resp, _ := http.Get("https://github.com/trending")
	tokenizer := html.NewTokenizer(resp.Body)
	for {
		tt := tokenizer.Next()
		t := tokenizer.Token()
		err := tokenizer.Err()
		if err == io.EOF {
			break
		}

		switch tt {
		case html.ErrorToken:
			log.Fatal(err)
		case html.StartTagToken:
			for _, a := range t.Attr {
				if a.Key == "href" && strings.Contains(a.Val, "/stargazers") {
					results = append(results, "https://github.com"+strings.ReplaceAll(a.Val, "/stargazers", ""))
				}
			}
		}
	}

	return results
}

func getAbout(url string) string {
	resp, _ := http.Get(url)
	tokenizer := html.NewTokenizer(resp.Body)
	for {
		tt := tokenizer.Next()
		t := tokenizer.Token()
		err := tokenizer.Err()
		if err == io.EOF {
			break
		}

		switch tt {
		case html.ErrorToken:
			log.Fatal(err)
		case html.StartTagToken:
			for _, a := range t.Attr {
				if a.Key == "class" && a.Val == "f4 mt-3" {
					tokenType := tokenizer.Next()
					if tokenType == html.TextToken {
						return strings.TrimSpace(tokenizer.Token().Data)
					}
				}
			}
		}
	}
	return "Description not available"
}
