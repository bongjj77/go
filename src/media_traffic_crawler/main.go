/*
Media traffic crawler
bongjj@gmail.com
20200301
- 구간 별 Streamer Http API를 호출 하여 traffic 정보 수집(json)
- 시차/트래픽 지연 정보 분석
- 정보 요청 출력(http/json)
- url_list.txt 정보 제공 Streamer url 목록(스트림 구간 순서로 설정)
*/
package main

import (
	"analyze"
	"bufio"
	"chart"
	"crawling"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	programName    = "Media traffic crawler"
	programVersion = "1.0"
)

// loadUrls : load url list file
func loadUrls(filePath string) (bool, []string) {
	file, error := os.Open(filePath)
	if error != nil {
		fmt.Println(filePath, "file open fail")
		return false, nil
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	urls := []string{}

	for {
		line, isPrefix, error := reader.ReadLine()
		if isPrefix == true || error != nil {
			break
		}
		urls = append(urls, string(line))
	}

	// urls print
	fmt.Println("Url count", len(urls))
	for index, url := range urls {
		fmt.Println(index, url)
	}

	return len(urls) != 0, urls
}

// readLocalHTML : local html file read
func readLocalHTML(path string) []byte {

	data, error := ioutil.ReadFile(path)
	if error != nil {
		fmt.Println(path, "read fail")
		return nil
	}

	return data
}

//====================================================================================================
// start
// - cpu count process
// - url list load
// - crawling timer
// - data analyze
// - http request process
//====================================================================================================
func main() {
	fmt.Println(programName, programVersion, time.Now().Format(time.RFC3339), "start")

	// max core use
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("Cpu core :", runtime.GOMAXPROCS(0))

	result, urls := loadUrls("./url_list.txt")
	if result == false {
		fmt.Println("Url list load fail")
		return
	}

	// key : first traffic StreamName
	collectCount := len(urls)
	collectedData := chart.NewCollectedData(len(urls))
	dataMutex := new(sync.Mutex)

	// crarwing time loop - go routine
	ticker := time.NewTicker(time.Millisecond * 5000)
	go func() {
		for start := range ticker.C {

			// crawling
			traffics := crawling.Crawling(urls)
			if len(traffics) == collectCount || len(traffics[0].StreamList) == 0 {
				fmt.Println("Traffics count fail")
				continue
			}

			// traffic analize
			analyzeData := analyze.Analyze(traffics)
			fmt.Println(analyzeData)

			// ------------------------------ sync start ------------------------------
			dataMutex.Lock()

			// time append
			collectedData.Times = append(collectedData.Times, start)

			if len(collectedData.Stream) == 0 {
				collectedData.Stream = traffics[0].StreamList[0].Stream
			}

			// latency append per section
			for _, latency := range analyzeData.LatencyList {
				collectedData.Sections[latency.Section].LatencyList = append(collectedData.Sections[latency.Section].LatencyList, latency.Latency)

				if len(collectedData.Sections[latency.Section].Host) == 0 {
					collectedData.Sections[latency.Section].Host = latency.Host
				}
			}

			//  ------------------------------ sync end ------------------------------
			dataMutex.Unlock()

			// process duration
			fmt.Println("Crawing :", start.Format(time.RFC3339), "duration :", time.Now().Sub(start).Milliseconds())

		}
	}()

	// file handle
	http.Handle("/http_file/", http.StripPrefix("/http_file/", http.FileServer(http.Dir("http_file"))))

	// root
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(programName + programVersion))
	})

	// get_latency_data.api
	http.HandleFunc("/get_latency_data.api", func(writer http.ResponseWriter, request *http.Request) {

		// ------------------------------ sync start ------------------------------
		dataMutex.Lock()

		chartHTML := chart.MakeLatecyChart(collectedData)
		fmt.Println("chart html :", chartHTML)
		//  ------------------------------ sync end ------------------------------
		dataMutex.Unlock()

		writer.Header().Set("Content-Type", "text/html") // HTML 헤더 설정
		writer.Write([]byte(chartHTML))
	})

	if error := http.ListenAndServe(":8080", nil); error != nil {
		fmt.Println("http server fail")
		return
	}

	// key input wait
	fmt.Scanln()

	// time loop close
	ticker.Stop()

	fmt.Println(programName, programVersion, time.Now().Format(time.RFC3339), "end")
}
