package main

import (
	"fmt"
	"runtime"
	"time"
)

const (
	programName    = "Media traffic crawler"
	programVersion = "1.0"
)

func crawing(urls []string) {

	// urls print
	for index, url := range urls {
		fmt.Println(index, url)
	}
	urlCount := len(urls)
	fmt.Println("Url count", urlCount)

	// timer create

	//

}

func main() {

	fmt.Println(programName, programVersion, time.Now().Format(time.RFC3339), "start")

	// max core use
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("Cpu core :", runtime.GOMAXPROCS(0))

	urls := []string{
		"http://wwww.naver.com",
		"http://www.google.com",
	}

	crawing(urls)

	// key input wait
	//fmt.Scanln()
	fmt.Println(programName, programVersion, time.Now().Format(time.RFC3339), "end")

}
