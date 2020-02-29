package analyze_test

import (
	. "analyze"
	. "crawling"
	. "testing"
)

// TestAnalyze : urls crawling
func TestAnalyze(t *T) {

	traffics := []*Traffic{
		NewTraffic("www.google.co.kr", 0, []byte(TestTrafficJSON)),
		NewTraffic("www.google.co.kr", 0, []byte(TestTrafficJSON)),
		NewTraffic("www.google.co.kr", 0, []byte(TestTrafficJSON))}

	Analyze(traffics)

}
