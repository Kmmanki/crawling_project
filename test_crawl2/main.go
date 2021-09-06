package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"test_crawl2/DTO"
	"test_crawl2/utils"

	"time"

	"github.com/PuerkitoBio/goquery"
)

var ApiBaseURL = "https://openapi.naver.com/v1/search/blog?display=100&sort=date&query="
var start = 1
var keyword = url.QueryEscape("삼성 TV")
var config utils.PrivateConfig
var err error

var blogBaseURL = "https://blog.naver.com"

// var config utils.PrivateConfig

func main() {
	startTime := time.Now()

	config, err = utils.LoadConfig()
	utils.ErrChecker(err)

	targetList := CallNaverAPI(ApiBaseURL, keyword)
	log.Println("api 추출 개수: ", len(targetList))

	ch := make(chan DTO.Scrap_result)
	for index, item := range targetList {
		if index%300 == 0 {
			time.Sleep(5 * time.Millisecond)
		}
		go crawlBlog(item, ch)
	}

	putCount := 0
	for i := 0; i < len(targetList); i++ {
		result := <-ch

		var re = regexp.MustCompile(`["']`)                      //따옴표가 존재 할 시 Joson 구조 만드는데 문제가 생김 추후 html 태그 등을 처리하는 정규식이 필요해 보임
		result.Content = re.ReplaceAllString(result.Content, ``) //Content에서 정규식 사용
		result.Title = re.ReplaceAllString(result.Title, ``)     // Title에서 정규식 사용

		putCount += utils.EsPut(result)
	}

	endTime := time.Now()
	utils.Insert_history(startTime, endTime, putCount, len(targetList))
}

func crawlBlog(target DTO.NaverBlogApiItem, ch chan<- DTO.Scrap_result) {
	url := target.Link
	if !strings.Contains(url, "blog.naver") {
		return
	}
	iframeSrc := getIframeURL(url)
	resp, err := http.Get(blogBaseURL + iframeSrc)
	utils.ErrChecker(err)
	statusCodeChecker(resp, "크롤링 블로그", url)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	defer resp.Body.Close()
	utils.ErrChecker(err)

	content := doc.Find(".se-main-container  .se-text-paragraph").Text() // Document에서 Content에 관련된 Class를 찾아내고 Text를 로드

	// Title과 PostDate를 사용한 MD5해쉬를 사용, 추후 엘라스틱에 있는 글이라면 넣지 않는 방향으로 진행 해야함.
	scrapId := utils.GetMD5Hash(target.Title, target.CustomPostDate.Format("2006-01-02 15:04:05"))

	//채널을 통해 Scrap_result 구조체를 반환
	ch <- DTO.Scrap_result{Title: target.Title, Link: target.Link, Content: content,
		CustomPostDate: target.CustomPostDate,
		ScrapDate:      time.Now(), ScrapId: scrapId}
}

//네이버 블로그는 iframe에 쌓인 형태로 구성되어 있음, 그래서
//첫째 블로그 url을 가져오고
//둘째 블로그 내용에서 iframe url을 가져온 뒤
//셋째 다시 iframe의 url을 호출하여 내용을 추출해야함
func getIframeURL(url string) string {
	var iframeURL string

	resp, err := http.Get(url)
	utils.ErrChecker(err)
	statusCodeChecker(resp, "아이프레임", url)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	defer resp.Body.Close()
	utils.ErrChecker(err)

	doc.Find("iframe").Each(func(i int, s *goquery.Selection) {
		iframeURL, _ = s.Attr("src")
	})
	return strings.Split(iframeURL, " ")[0]
}

func CallNaverAPI(apiURL string, keyword string) []DTO.NaverBlogApiItem {
	log.Println("api call start")

	var targetList []DTO.NaverBlogApiItem
	onedaysAgo := time.Now().AddDate(0, 0, -1) //.Format("2006-01-02")
	var bloglist DTO.NaverBlogApiStruct
	isBreak := false

	for {
		client := &http.Client{}

		url := apiURL + keyword + "&start=" + strconv.Itoa(start)
		req, err := http.NewRequest("GET", url, nil)
		req.Header.Set("X-Naver-Client-Id", config.NaverAPI.ClientCode)
		req.Header.Set("X-Naver-Client-Secret", config.NaverAPI.SecretCode)

		resp, err := client.Do(req)
		utils.ErrChecker(err)
		fmt.Println(url)
		statusCodeChecker(resp, "네이버 api호출", url)

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		// body := string(bodyBytes)
		decoder := json.NewDecoder(bytes.NewReader(bodyBytes))
		err = decoder.Decode(&bloglist)
		utils.ErrChecker(err)
		// fmt.Println("3일전 날짜", threeDaysAgo.Format("2006-01-02"))
		for i := 0; i < len(bloglist.Items); i++ {
			strPostDate := bloglist.Items[i].Postdate

			y, _ := strconv.Atoi(strPostDate[0:4])
			m, _ := strconv.Atoi(strPostDate[4:6])
			d, _ := strconv.Atoi(strPostDate[6:])
			postDate := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)

			bloglist.Items[i].CustomPostDate = postDate

			if postDate.Before(onedaysAgo) {
				isBreak = true
			}
			targetList = append(targetList, bloglist.Items[i])
		}
		time.Sleep(time.Millisecond * 5)
		fmt.Println(bloglist.Items[0].CustomPostDate.Format("2006-01-02"))
		if isBreak {
			break
		}
		start += 100
	}
	log.Println("api call end")

	return targetList

}

func statusCodeChecker(resp *http.Response, target string, url string) {
	code := resp.StatusCode
	if code != 200 {
		log.Fatalln("StatusCode: ", code, "target :"+target, "url: ", url)
	}
}
