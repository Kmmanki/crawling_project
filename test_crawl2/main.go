package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"test_crawl2/DTO"
	"test_crawl2/utils"

	"time"

	"github.com/PuerkitoBio/goquery"
)

var ApiBaseURL = "https://openapi.naver.com/v1/search/blog?display=3&sort=date&query="
var start = 0
var keyword = url.QueryEscape("삼성 TV")
var startPage = 1
var config utils.PrivateConfig
var err error

var blogBaseURL = "https://blog.naver.com"

// var config utils.PrivateConfig

func main() {

	config, err = utils.LoadConfig()
	utils.ErrChecker(err)

	targetList := CallNaverAPI(ApiBaseURL, keyword)

	crawlBlog(targetList[0])

}

func crawlBlog(target DTO.NaverBlogApiItem) {
	println(target.Link)
	url := target.Link
	iframeSrc := getIframeURL(url)

	resp, err := http.Get(blogBaseURL + iframeSrc)
	utils.ErrChecker(err)
	statusCodeChecker(resp)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	defer resp.Body.Close()
	utils.ErrChecker(err)

	doc.Find(".se-main-container  .se-text-paragraph").Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Text())
	})

}

func getIframeURL(url string) string {
	var iframeURL string
	resp, err := http.Get(url)
	utils.ErrChecker(err)
	statusCodeChecker(resp)

	// bodyBytes, err := ioutil.ReadAll(resp.Body)
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	defer resp.Body.Close()
	utils.ErrChecker(err)

	doc.Find("iframe").Each(func(i int, s *goquery.Selection) {
		iframeURL, _ = s.Attr("src")
	})
	return strings.Split(iframeURL, " ")[0]
}

func CallNaverAPI(url string, keyword string) []DTO.NaverBlogApiItem {
	var targetList []DTO.NaverBlogApiItem
	threeDaysAgo := time.Now().AddDate(0, 0, -3) //.Format("2006-01-02")
	var bloglist DTO.NaverBlogApiStruct
	isBreak := false

	for {
		client := &http.Client{}

		url = url + keyword + "&start=" + strconv.Itoa(startPage)
		req, err := http.NewRequest("GET", url, nil)
		req.Header.Set("X-Naver-Client-Id", config.NaverAPI.ClientCode)
		req.Header.Set("X-Naver-Client-Secret", config.NaverAPI.SecretCode)

		resp, err := client.Do(req)
		utils.ErrChecker(err)
		statusCodeChecker(resp)

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		// body := string(bodyBytes)
		decoder := json.NewDecoder(bytes.NewReader(bodyBytes))
		err = decoder.Decode(&bloglist)
		utils.ErrChecker(err)

		fmt.Println("3일전 날짜", threeDaysAgo.Format("2006-01-02"))
		for i := 0; i < len(bloglist.Items); i++ {
			strPostDate := bloglist.Items[i].Postdate

			y, _ := strconv.Atoi(strPostDate[0:4])
			m, _ := strconv.Atoi(strPostDate[4:6])
			d, _ := strconv.Atoi(strPostDate[6:])
			postDate := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)

			bloglist.Items[i].CustomPostDate = postDate

			if postDate.Before(threeDaysAgo) {
				isBreak = true
			}
			targetList = append(targetList, bloglist.Items[i])
		}

		if isBreak || startPage > 2 {
			break
		}
		startPage += 1
	}

	// doc, err := goquery.NewDocumentFromReader(resp.Body)

	// errChecker(err)
	// doc.Find(".\"total_sub\"").Each(func(i int, s *goquery.Selection) {
	// 	fmt.Println(s.Length())
	// })
	return targetList

}

func statusCodeChecker(resp *http.Response) {
	code := resp.StatusCode
	if code != 200 {
		log.Fatalln("StatusCode: ", code)
	}
}
