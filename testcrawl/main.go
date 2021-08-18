package main

import (
	"fmt"
	"net/http"
)

type requestResult struct {
	url        string
	status     string
	statusCode string
}

func main() {
	c := make(chan requestResult)
	results := make(map[string]string)
	urls := []string{
		"https://www.airbnb.co.kr/",
		"https://www.google.com/",
		"https://www.amazon.com/",
		"https://www.reddit.com/",
		"https://www.soundcloud.com/",
		"https://www.facebook.com/",
		"https://www.instagram.com/",
	}

	for _, url := range urls {
		go hitURL(url, c)

	}

	for i := 0; i < len(urls); i++ {
		result := <-c
		results[result.url] = result.status
	}
	fmt.Println(results)

}

func hitURL(url string, c chan<- requestResult) {
	fmt.Println("URL Checking...")

	resp, err := http.Get(url)
	status := "Ok"
	if err != nil || resp.StatusCode >= 400 {
		status = "Faild"
		c <- requestResult{url: url, status: status, statusCode: resp.Status}
	} else {
		c <- requestResult{url: url, status: status, statusCode: resp.Status}

	}

}
