package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/test_crawl2/utils"
)

var baseURL = "https://openapi.naver.com/v1/search/blog?display=3&sort=date&query="
var clientCode = ""
var secretCode = ""
var start = 0
var keyword = url.QueryEscape("삼성 TV")
var startPage = 1

var config string

func main() {

	config = utils.LoadConfig
	println(config)

	//hitURL(baseURL, keyword)
}

func hitURL(url string, keyword string) {

	url = url + keyword + "&start=" + strconv.Itoa(startPage)

	client := &http.Client{}
	fmt.Println("checking: ", url)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Naver-Client-Id", clientCode)
	req.Header.Set("X-Naver-Client-Secret", secretCode)

	resp, err := client.Do(req)
	errChecker(err)
	statusCodeChecker(resp)

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	print("bodt:", string(bodyBytes))
	// doc, err := goquery.NewDocumentFromReader(resp.Body)

	// errChecker(err)
	// doc.Find(".\"total_sub\"").Each(func(i int, s *goquery.Selection) {
	// 	fmt.Println(s.Length())
	// })

}

func errChecker(err error) {
	if err != nil {
		log.Fatalln("err: ", err)
	}
}

func statusCodeChecker(resp *http.Response) {
	code := resp.StatusCode
	if code != 200 {
		log.Fatalln("StatusCode: ", code)
	}
}
