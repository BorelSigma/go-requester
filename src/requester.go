package main

import (
"fmt"
"os"
"net/http"
"time"
"strconv"
)

func MakeRequest(url string, ch chan<-string) {
	start := time.Now()
	resp, err := http.Post(url, "", nil)

	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2f elapsed with response length: %d %s", secs, url, err, resp)
}

func main() {
	url := os.Args[1]
	start := time.Now()
	ch := make(chan string)
	for _,seconds := range os.Args[2:]{
		i, _ := strconv.Atoi(seconds)
		time.Sleep(time.Duration(i) * time.Millisecond)
		go MakeRequest(url, ch)
	}

	for range os.Args[1:]{
		fmt.Println(<-ch)
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}