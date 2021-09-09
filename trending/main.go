package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"strings"
  "encoding/json"
  "bytes"
  "os"
)

type Message struct {
	Project string
	About   string
}

type Payload struct {
  Username string `json:"username"`
  Avatar_url string `json:"avatar_url"`
  Content string `json:"content"`
  Embeds [1]Embed `json:"embeds"`
}

type Embed struct {
  Author Author `json:"author"`
  Url string `json:"url"`
  Title string `json:"title"`
  Description string `json:"description"`
  Color int `json:"color"`
  Fields []Field `json:"fields"`
  Footer Footer `json:"footer"`
}

type Footer struct {
  Text string `json:"text"`
  Icon_url string `json:"icon_url"`
}

type Author struct {
  Name string `json:"name"`
  Url string `json:"url"`
  Icon_url string `json:"icon_url"`
}

type Field struct {
  Name string `json:"name"`
  Value string `json:"value"`
  Inline string `json:"inline"`
}

func main() {
	c := make(chan Message)

	urls := getTrending()

	for _, url := range urls {
		go getOutput(url, c)
	}

	result := make([]Message, len(urls))
  f := make([]Field, len(urls))
	for i, _ := range result {
		result[i] = <-c
    f[i] = Field{
      Name: result[i].Project,
      Value: result[i].About,
      Inline: "false",
    }
  }

  emb := Embed{
    Author: Author{
      Name: "Trending Bot",
      Url: "https://google.com",
      Icon_url: "https://i.imgur.com/R66g1Pe.jp",
    },
    Url: "https://github.com",
    Title: "Trending Topics",
    Description: "Automated messages",
    Color: 15258703,
    Fields: f,
    Footer: Footer{
      Text: "Brought to you by cachemeoutside",
      Icon_url: "https://i.imgur.com/R66g1Pe.jpg",
    },
  }

  payload := Payload{
    Username: "Trending Bot",
    Avatar_url: "https://i.imgur.com/R66g1Pe.jpg",
    Content: "Content",
    Embeds: [1]Embed{emb},
  }

  b, err := json.Marshal(payload)
  if err != nil {
    fmt.Println(err)
    return
  }

  resp, err := http.Post(os.Getenv("DISCORD_TRENDING_WEBHOOK"), "application/json", bytes.NewBuffer(b))

  if err != nil {
    fmt.Println(err)
    return
  }

  fmt.Println(resp)
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
            v := strings.TrimSpace(tokenizer.Token().Data)
            if v != "" {
  						msg.About = v
            }
					}
				}
			}
		}
	}
	c <- msg
}
