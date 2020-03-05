package crawling_test

import (
	. "crawling"
	. "fmt"
	. "testing"
)

// TestReadTraffic : restreamer output json data parse
func TestReadTraffic(t *T) {

	result := make(chan *Traffic)

	go ReadTraffic("http://www.google.com", 0, result)

	Println(<-result)

}

// TestCrawling : urls crawling
// "TODO : single" test 주석 확인/제거
func TestCrawling(t *T) {

	urls := make([]string, 3)
	urls[0] = "http://www.google.com"
	urls[1] = "http://www.google.com"
	urls[2] = "http://www.google.com"

	traffics := Crawling(urls)

	for _, traffic := range traffics {
		Println(traffic)
	}
}
