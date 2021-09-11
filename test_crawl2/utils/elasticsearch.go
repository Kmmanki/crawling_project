package utils

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"test_crawl2/DTO"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"golang.org/x/net/html/charset"
)

func EsPut(result DTO.Scrap_result) int {
	es, err := elasticsearch.NewDefaultClient()
	// Build the request body.
	var b strings.Builder
	if len(result.Content) >= 10000 {
		result.Content = result.Content[0:10000]
	}

	b.WriteString(`{`)
	b.WriteString(`"Title" : "`)
	b.WriteString(result.Title)
	b.WriteString(`",`)

	b.WriteString(`"Content" : "`)
	b.WriteString(result.Content)
	b.WriteString(`",`)

	b.WriteString(`"Link" : "`)
	b.WriteString(result.Link)
	b.WriteString(`",`)

	b.WriteString(`"PostDate" : "`)
	b.WriteString(result.CustomPostDate.Format("2006-01-02 15:04:05"))
	b.WriteString(`",`)

	b.WriteString(`"ScrapDate" : "`)
	b.WriteString(result.ScrapDate.Format("2006-01-02 15:04:05"))
	b.WriteString(`",`)

	b.WriteString(`"ScrapId" : "`)
	b.WriteString(result.ScrapId)
	b.WriteString(`"}`)

	// Set up the request object.
	body, _ := charset.NewReader(strings.NewReader(b.String()), "UTF-8")
	req := esapi.IndexRequest{
		Index:      "naver_scraper",
		DocumentID: result.ScrapId,
		Body:       body,
		Refresh:    "true",
	}
	// Perform the request with the client.
	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		res.Body.Close()

	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[%s] Error indexing document", res.Status())
		log.Println("url: ", result.Link)
		// fmt.Println(b.String())
		// log.Println("ReqBody: ", strings.NewReader(b.String()))
		// log.Fatal(len(result.Content))
		log.Fatalln("RespBody: ", res.String())

		// log.Println("RespBody: ", strings.NewReader(res.Body))

		return 0
	} else {
		// Deserialize the response into a map.
		// log.Fatalln("sssss", strings.NewReader(b.String()))
		// fmt.Println(b.String())
		// log.Fatalln()

		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and indexed document version.
			log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
		}
		return 1
	}
}
