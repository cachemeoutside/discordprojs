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
					fmt.Printf("https://github.com%s\n", strings.ReplaceAll(a.Val, "/stargazers", ""))
				}
			}
		}
	}
}
