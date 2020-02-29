/*
Media traffic crawler
bongjj@gmail.com
20200301
동일 스트림의 구간 별 Streamer demone API를 동시 호출 하여 json 형식의 정보 수집 하여 시차/트래픽 정보를 처리
url_list.txt 정보 제공 url 목록
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
			analyze.Analyze(traffics)

			// process duration
			fmt.Println("Crawing :", start.Format(time.RFC3339), "duration :", time.Now().Sub(start).Milliseconds())

		}
	}()

	// key input wait
	fmt.Scanln()

	// time loop close
	ticker.Stop()

	fmt.Println(programName, programVersion, time.Now().Format(time.RFC3339), "end")
}
