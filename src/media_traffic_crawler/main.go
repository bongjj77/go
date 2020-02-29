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
	"os"
	"runtime"
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

	// crarwing time loop - go routine
	ticker := time.NewTicker(time.Millisecond * 5000)
	go func() {
		for start := range ticker.C {

			// crawling
			traffics := crawling.Crawling(urls)

			// traffic analize
			streamAnalyze := analyze.Analyze(traffics)

			fmt.Println(streamAnalyze)

			// TODO : data save

			// process duration
			fmt.Println("Crawing :", start.Format(time.RFC3339), "duration :", time.Now().Sub(start).Milliseconds())

		}
	}()

	// TODO : http server
	// - 확인 요청 시간 기준 최대 개수 설정 하여 전송
	// - 그래프 출력

	// key input wait
	fmt.Scanln()

	// time loop close
	ticker.Stop()

	fmt.Println(programName, programVersion, time.Now().Format(time.RFC3339), "end")
}
