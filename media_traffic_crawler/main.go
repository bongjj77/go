/*
Media traffic crawler
bongjj@gmail.com
20200301
restream demone 에서 json 형식의 정보를 특정 시간에 동시 수집 하여  시차/트래픽 정보를 처리
url_list.txt 정보 제공 url 목록
*/
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"time"
)

const (
	programName    = "Media traffic crawler"
	programVersion = "1.0"
)

type people struct {
	Number int `json:"number"`
}

//====================================================================================================
// Http response data parse
// - json data
// - single input
// - multi output
// - timestamp/bitrate/fps
// ex)
//====================================================================================================
func dataParse(data []byte) string {

	people1 := people{}
	error := json.Unmarshal(data, &people1)
	if error != nil {
		fmt.Println("Json parae fail")
		return ""
	}

	fmt.Println(people1.Number)
	return ""
}

//====================================================================================================
// Read url(restream demon api)
//====================================================================================================
func readURL(url string, result chan<- string) {

	response, error := http.Get(url)
	if error != nil {
		result <- (url + " http get fail")
		return
	}
	defer response.Body.Close()

	data, error := ioutil.ReadAll(response.Body)
	if error != nil {
		result <- (url + " body read fail")
		return
	}

	dataParse(data)
	fmt.Println(string(data))

	result <- (url + " read comleted")
}

//====================================================================================================
// Crawing
//====================================================================================================
func crawling(urls []string) {

	// urls print
	for index, url := range urls {
		fmt.Println(index, url)
	}
	urlCount := len(urls)
	fmt.Println("Url count", urlCount)

	result := make(chan string)

	// read go rutin
	for index := 0; index < urlCount; index++ {
		go readURL(urls[index], result)
	}

	// complted wait
	for index := 0; index < urlCount; index++ {
		fmt.Println(<-result)
	}

}

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

	// time loop
	ticker := time.NewTicker(time.Millisecond * 5000)
	go func() {
		for start := range ticker.C {
			crawling(urls)
			fmt.Println("Crawing :", start.Format(time.RFC3339), "duration : ", time.Now().Sub(start).Milliseconds())
		}
	}()

	// key input wait
	fmt.Scanln()

	// time loop close
	ticker.Stop()

	fmt.Println(programName, programVersion, time.Now().Format(time.RFC3339), "end")
}
