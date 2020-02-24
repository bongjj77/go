package trafficparse

import (
	"encoding/json"
	"fmt"
)

type people struct {
	Number int `json:"number"`
}

// RestreamerDataParse : restreamer output json data parse
func RestreamerDataParse(data []byte) string {

	people1 := people{}
	error := json.Unmarshal(data, &people1)
	if error != nil {
		fmt.Println("Json parae fail")
		return ""
	}

	fmt.Println(people1.Number)
	return ""
}
