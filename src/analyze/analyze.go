package analyze

import (
	"crawling"
	"fmt"
)

// Analyze :  Traffic struct analyze
func Analyze(traffics []*crawling.Traffic) {

	// print
	for _, traffic := range traffics {
		fmt.Println(traffic)
	}

}
