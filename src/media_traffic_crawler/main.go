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
	"crawling"
	"fmt"
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

//====================================================================================================
// load url list file
//====================================================================================================
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

type collectData struct {
	stream      string
	analyzeList []*analyze.StreamAnalyze
}

//====================================================================================================
// start
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
	collectDataMap := make(map[string]*collectData)
	dataMutex := new(sync.Mutex)

	// crarwing time loop - go routine
	ticker := time.NewTicker(time.Millisecond * 5000)
	go func() {
		for start := range ticker.C {

			// crawling
			traffics := crawling.Crawling(urls)
			if len(traffics) == 0 || len(traffics[0].StreamList) == 0 {
				fmt.Println("Traffics size zero")
				continue
			}

			stream := traffics[0].StreamList[0].Stream

			// traffic analize
			streamAnalyze := analyze.Analyze(traffics)

			fmt.Println(streamAnalyze)

			// ------------------------------ sync start ------------------------------
			dataMutex.Lock()

			// data insert
			if _, exist := collectDataMap[stream]; exist != true {
				collectDataMap[stream] = &collectData{stream, make([]*analyze.StreamAnalyze, 0)}
			}
			collectDataMap[stream].analyzeList = append(collectDataMap[stream].analyzeList, streamAnalyze)

			//  ------------------------------ sync end ------------------------------
			dataMutex.Unlock()

			// process duration
			fmt.Println("Crawing :", start.Format(time.RFC3339), "duration :", time.Now().Sub(start).Milliseconds())

		}
	}()

	// TODO : http server
	// - 확인 요청 시간 기준 최대 개수 설정 하여 전송
	// - 그래프 출력
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		html := `
				<html>
				<head>
					<script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/Chart.js/2.8.0/Chart.bundle.min.js"></script>
					<script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/Chart.js/2.8.0/Chart.min.js"></script>
				</head>
				<canvas id="myChart" width="400" height="400"></canvas>
				<script>
					var ctx = document.getElementById('myChart');
					var myChart = new Chart(ctx, {
						type: 'line',
						data: {
							labels: ['Red', 'Blue', 'Yellow', 'Green', 'Purple', 'Orange'],
							datasets: [{
								label: '# of Votes',
								data: [12, 19, 3, 5, 2, 3],
								backgroundColor: [
									'rgba(255, 99, 132, 0.2)',
									'rgba(54, 162, 235, 0.2)',
									'rgba(255, 206, 86, 0.2)',
									'rgba(75, 192, 192, 0.2)',
									'rgba(153, 102, 255, 0.2)',
									'rgba(255, 159, 64, 0.2)'
								],
								borderColor: [
									'rgba(255, 99, 132, 1)',
									'rgba(54, 162, 235, 1)',
									'rgba(255, 206, 86, 1)',
									'rgba(75, 192, 192, 1)',
									'rgba(153, 102, 255, 1)',
									'rgba(255, 159, 64, 1)'
								],
								borderWidth: 1
							}]
						},
						options: {
							responsive: false,
							scales: {
								yAxes: [{
									ticks: {
										beginAtZero: true
									}
								}]
							},
						}
					});
				</script>
				</html>
				`

		writer.Header().Set("Content-Type", "text/html") // HTML 헤더 설정
		writer.Write([]byte(html))                       // 웹 브라우저에 응답
	})

	http.Handle("/http_file/", http.StripPrefix("/http_file/", http.FileServer(http.Dir("http_file"))))

	http.HandleFunc("/get_latency_data.api", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("get lantency test test"))
	})

	http.ListenAndServe(":8080", nil)

	// key input wait
	fmt.Scanln()

	// time loop close
	ticker.Stop()

	fmt.Println(programName, programVersion, time.Now().Format(time.RFC3339), "end")
}
