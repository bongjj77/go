package trafficparse_test

import (
	"fmt"
	"io/ioutil"
	"testing"
	"trafficparse"
)

// TestStreamerDataParse : restreamer output json data parse
func TestStreamerDataParse(t *testing.T) {
	data, error := ioutil.ReadFile("./traffic.json")
	if error != nil {
		t.Error("Json file open fail")
		return
	}

	result, parseData := trafficparse.StreamerDataParse(data)
	if result == false || parseData == nil {
		t.Error("Data parse fail")
		return
	}

	fmt.Println(parseData)

}

// TestReadTraffic : restreamer output json data parse
func TestReadTraffic(t *testing.T) {

	result := make(chan *trafficparse.Traffic)

	go trafficparse.ReadTraffic("http://www.google.com", 0, result)

	fmt.Println(<-result)

}
