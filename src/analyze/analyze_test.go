package analyze_test

import (
	. "analyze"
	. "crawling"
	"fmt"
	. "testing"
)

// TestAnalyze : urls crawling
// "TODO : single" test 주석 확인/제거
func TestAnalyze(t *T) {

	traffics := []*Traffic{
		NewTraffic("www.google.co.kr", 0, []byte(TestTrafficJSON)),
		NewTraffic("www.google1.co.kr", 1, []byte(TestTrafficJSON1)),
		NewTraffic("www.google2.co.kr", 2, []byte(TestTrafficJSON2))}

	streamAnalyze := Analyze(traffics)

	fmt.Print(streamAnalyze)
}
