package trafficparse_test

import (
	. "fmt"
	. "testing"
	. "trafficparse"
)

// TestReadTraffic : restreamer output json data parse
func TestReadTraffic(t *T) {

	result := make(chan *Traffic)

	go ReadTraffic("http://www.google.com", 0, result)

	Println(<-result)

}
