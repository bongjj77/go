package trafficparse_test

import (
	"fmt"
	"io/ioutil"
	"testing"
	"trafficparse"
)

// StreamerDataParse : restreamer output json data parse
func TestStreamerDataParse(t *testing.T) {

	t.Error("Json file open fail")

	// temp test
	data, error := ioutil.ReadFile("./traffic.json")
	if error != nil {
		t.Error("Json file open fail")
		return
	}

	result, parseData := trafficparse.StreamerDataParse(data)
	if result == false {
		t.Error("Data parse fail")
		return
	}

	fmt.Println(parseData)

}
