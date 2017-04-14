package main

import (
"fmt"
"os"
"net/http"
"time"
"strconv"
"bytes"
"io/ioutil"
"bufio"
"encoding/json"
)

type Request struct {
	Variables map[string] interface{}  `json:"variables"`
}

func MakeRequest( id int, url string, ch chan<-string) {
	start := time.Now()
	var requestJson Request
	lines := File2lines(strconv.Itoa(id)+".token")
	vars := make(map[string] interface{})
	vars["output"] = "/opt/app/results/"
	vars["music_data_type"] = "audio/midi"
	vars["input"] = "/opt/app/music_files/"
	vars["tmp_dir"] = "/opt/app/tmp/"
	vars["file_name"] = "Eminem_-_The_Real_Slim_Shady.mid"
	vars["channels_value"] = 1
	vars["normalize_value"] = 0.6
	vars["limiter_value"] = 0.4
	vars["fading_value"] = 0.13
	vars["sample_size_value"] = 8
	vars["sample_rate_value"] = 4300
	vars["normalize_webservice_url"] = lines[0]
	vars["limiter_webservice_url"] = lines[1]
	vars["sample_rate_webservice_url"] = lines[2]
	vars["channels_webservice_url"] = lines[3]
	vars["fading_webservice_url"] = lines[4]
	vars["sample_size_webservice_url"] = lines[5]
    
    requestJson.Variables = vars
	jsonStr, err := json.Marshal(requestJson)
	if err != nil {
		fmt.Println("error:", err)
	}
	
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2f elapsed", secs)
}

func File2lines(filePath string) []string {
      f, err := os.Open(filePath)
      if err != nil {
              panic(err)
      }
      defer f.Close()

      var lines []string
      scanner := bufio.NewScanner(f)
      for scanner.Scan() {
              lines = append(lines, scanner.Text())
      }
      if err := scanner.Err(); err != nil {
              fmt.Fprintln(os.Stderr, err)
      }

      return lines
}

func main() {
	url := os.Args[1]
	start := time.Now()
	ch := make(chan string)

	for id ,seconds := range os.Args[2:]{
		i, _ := strconv.Atoi(seconds)
		time.Sleep(time.Duration(i) * time.Millisecond)
		go MakeRequest(id, url, ch)
	}

	for range os.Args[2:]{
		fmt.Println(<-ch)
	}
	fmt.Printf("%.2fs total elapsed\n", time.Since(start).Seconds())
}