package analyze

import (
	"crawling"
	"fmt"
	"math/rand"
	"time"
)

// Latency : latency info
type Latency struct {
	Section int
	Host    string
	Latency int64 // millisecond
}

// StreamAnalyze : traffic analyze result
type StreamAnalyze struct {
	CreateTime  time.Time
	LatencyList []Latency
}

// NewStreamAnalyze : traffic create
func NewStreamAnalyze() *StreamAnalyze {
	streamAnalyze := new(StreamAnalyze)

	streamAnalyze.CreateTime = time.Now()
	streamAnalyze.LatencyList = make([]Latency, 0)
	return streamAnalyze
}

// Analyze :  Traffic struct analyze
func Analyze(traffics []*crawling.Traffic) *StreamAnalyze {

	// print
	streamAnalyze := NewStreamAnalyze()
	var timestamp uint64
	for index, traffic := range traffics {

		fmt.Println(traffic)

		if len(traffic.StreamList) == 0 {
			fmt.Println("Not exist stream")
			break
		}

		if index == 0 {
			timestamp = traffic.StreamList[0].Input.VideoTimestamp
		}

		// test
		timestamp = uint64(rand.Intn(100))
		streamAnalyze.LatencyList = append(streamAnalyze.LatencyList, Latency{traffic.SectionNumber, traffic.Host.HostName, int64(timestamp)})

		//streamAnalyze.LatencyList = append(streamAnalyze.LatencyList, Latency{traffic.SectionNumber, traffic.Host.HostName,, int64(timestamp - traffic.StreamList[0].Input.VideoTimestamp)})

	}

	return streamAnalyze
}
