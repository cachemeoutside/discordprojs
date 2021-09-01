package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"strings"
)

type Message struct {
	Project string
	About   string
}

func main() {
	c := make(chan Message)

	urls := getTrending()

	for _, url := range urls {
		go getOutput(url, c)
	}

	result := make([]Message, len(urls))
	for i, _ := range result {
		result[i] = <-c
		fmt.Printf("%s | %s\n", result[i].Project, result[i].About)
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

func getOutput(url string, c chan Message) {
	resp, _ := http.Get(url)
	tokenizer := html.NewTokenizer(resp.Body)
	msg := Message{url, "Description not available"}
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
						msg.About = strings.TrimSpace(tokenizer.Token().Data)
					}
				}
			}
		}
	}
	c <- msg
}
