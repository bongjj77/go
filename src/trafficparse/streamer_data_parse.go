package trafficparse

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

/*
 json to struct : https://mholt.github.io/json-to-go/

{
    "host": {
        "cpu": 0,
        "host_name": "test01",
        "input_count": 1,
        "ip": "127.0.0.1",
        "memory_total": 8127,
        "memory_used": 2161,
        "output_count": 1,
        "region": "korea",
        "time": "2020-02-24 17:00:28",
        "traffic": 927994458,
        "version": "1.0"
    },
    "stream_list":[{
            "stream": "app/123-456",
            "input": {
                "audio_framerate": 47.0,
                "audio_timestamp": 926229,
                "recv_traffic": 464221743,
                "remote": "127.0.0.1:6782",
                "send_traffic": 9638,
                "video_framerate": 30.0,
                "video_timestamp": 926200
            },
            "output": [
                {
                    "recv_traffic": 9787,
                    "remote": "127.0.0.1:1935",
                    "send_buffer": 0,
                    "send_traffic": 463753290,
                    "stream": "app/123-456"
                }
            ]
		}
	]
}
*/

// Traffic : json traffic struct
type Traffic struct {
	SectionNumber int

	Host struct {
		CPU         int    `json:"cpu"`
		HostName    string `json:"host_name"`
		InputCount  int    `json:"input_count"`
		IP          string `json:"ip"`
		MemoryTotal int    `json:"memory_total"`
		MemoryUsed  int    `json:"memory_used"`
		OutputCount int    `json:"output_count"`
		Region      string `json:"region"`
		Time        string `json:"time"`
		Traffic     int    `json:"traffic"`
		Version     string `json:"version"`
	} `json:"host"`
	StreamList []struct {
		Stream string `json:"stream"`
		Input  struct {
			AudioFramerate float64 `json:"audio_framerate"`
			AudioTimestamp int     `json:"audio_timestamp"`
			RecvTraffic    int     `json:"recv_traffic"`
			Remote         string  `json:"remote"`
			SendTraffic    int     `json:"send_traffic"`
			VideoFramerate float64 `json:"video_framerate"`
			VideoTimestamp int     `json:"video_timestamp"`
		} `json:"input"`
		Output []struct {
			RecvTraffic int    `json:"recv_traffic"`
			Remote      string `json:"remote"`
			SendBuffer  int    `json:"send_buffer"`
			SendTraffic int    `json:"send_traffic"`
			Stream      string `json:"stream"`
		} `json:"output"`
	} `json:"stream_list"`
}

// StreamerDataParse : restreamer output json data parse
func StreamerDataParse(data []byte) (bool, *Traffic) {

	traffic := &Traffic{}
	error := json.Unmarshal(data, traffic)
	if error != nil {
		fmt.Println("Json parae fail")
		return false, nil
	}

	return true, traffic
}

// ReadTraffic : Read url(restream demon api)
func ReadTraffic(url string, sectionNumber int, traffic chan<- *Traffic) {

	response, error := http.Get(url)
	if error != nil {
		fmt.Println(url, "http get fail")
		traffic <- nil
		return
	}
	defer response.Body.Close()

	data, error := ioutil.ReadAll(response.Body)
	if error != nil {
		fmt.Println(url, "body read fail")
		traffic <- nil
		return
	}

	tempData, error := ioutil.ReadFile("./traffic.json")
	if error != nil {
		fmt.Println(url, "json file read fail")
		traffic <- nil
		return
	}
	data = tempData

	result, parsedTraffic := StreamerDataParse(data)
	if result == false {
		fmt.Println(url, "data parse fail")
		traffic <- nil
		return
	}
	fmt.Println(url, "read comleted")

	parsedTraffic.SectionNumber = sectionNumber
	traffic <- parsedTraffic
}
