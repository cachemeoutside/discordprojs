package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	getJson(getRealUrl())
}
func getRealUrl() string {
	resp, _ := http.Get("https://www.microsoft.com/en-us/download/confirmation.aspx?id=56519")
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
				if a.Key == "href" && strings.Contains(a.Val, ".json") {
					return a.Val
				}
			}
		}
	}

	return ""
}

func getJson(url string) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Print(string(b))
}
