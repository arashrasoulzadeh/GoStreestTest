package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/schollz/progressbar"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

//strucure of the yaml file
type conf struct {
	Hits   int                    `yaml:"hits"`
	Route  string                 `yaml:"route"`
	Code   int                    `yaml:"code"`
	Method string                 `yaml:"method"`
	Body   map[string]interface{} `yaml:"body"`
}

//get the configuration
func (c *conf) getConf() *conf {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

//send http call
func SendHttpRequest(method string, url string, body string, wg *sync.WaitGroup) (*http.Response, error) {
	method = strings.ToUpper(method)
	switch method {
	case http.MethodGet:
		fallthrough
	case http.MethodConnect:
		fallthrough
	case http.MethodDelete:
		fallthrough
	case http.MethodHead:
		fallthrough
	case http.MethodOptions:
		fallthrough
	case http.MethodPatch:
		fallthrough
	case http.MethodPost:
		fallthrough
	case http.MethodPut:
		fallthrough
	case http.MethodTrace:
		reqBody := []byte(body)
		req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		return resp, err
	}

	fmt.Println("No valid method was provided, using default method GET:")
	resp, err := http.Get(url)
	return resp, err
}

//recover if error occured in MakeRequest
func MakeRequestRecovery(wg *sync.WaitGroup) {
	defer wg.Done()
	if r := recover(); r != nil {
		fmt.Println("Recovered in f", r)
	}
}

//make request handler
func MakeRequest(thread int, url string, method string, body string, ch chan<- string, id int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, f *os.File, ferr error) {
	defer MakeRequestRecovery(wg)
	start := time.Now()
	resp, err := SendHttpRequest(method, url, body, wg)
	duration := time.Since(start).Seconds()
	if err != nil {
		// handle the error, often:
		writeToLog(thread, id, resp, err, duration, f, ferr)
		bar.Add(1)
		return
	}
	writeToLog(thread, id, resp, err, duration, f, ferr)
	bar.Add(1)
}

//write to the log
func writeToLog(thread int, id int, response *http.Response, e error, duration float64, f *os.File, ferr error) {
	if e != nil {
	} else {
		if _, ferr = f.WriteString(fmt.Sprintf("%d,%d,%d,%f,%s\n", thread, id, response.StatusCode, duration, "0")); ferr != nil {
			panic(ferr)
		}
	}
}

//clear log for use
func clearLog() {
	message := []byte("thread,id,code,duration,error\n")
	err := ioutil.WriteFile("log.csv", message, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

//worker
func worker(mainWaitGroup *sync.WaitGroup, thread int) {
	defer mainWaitGroup.Done()
	f, err := os.OpenFile("log.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer func() {
		f.Close()
	}()
	var wg sync.WaitGroup
	var c conf
	c.getConf()
	_ = time.Now()
	ch := make(chan string)
	wg.Add(c.Hits)
	bar := progressbar.New(c.Hits)
	bar.RenderBlank()
	var body string
	if c.Body != nil {
		marshaled, err := json.Marshal(c.Body)
		if err != nil {
			panic(err)
		}
		body = string(marshaled)
	} else {
		marshaled, _ := json.Marshal(map[string]interface{}{})
		body = string(marshaled)
	}
	for i := 1; i <= c.Hits; i++ {
		go MakeRequest(thread, c.Route, c.Method, body, ch, i, &wg, bar, f, err)
	}
	wg.Wait()
}

func single() {
	var wg sync.WaitGroup
	wg.Add(1)
	worker(&wg, 1)
}
func multiple() {
	input := os.Args
	if len(input) == 3 {
		counts := input[2]
		count, err := strconv.Atoi(counts)
		if err != nil {
			help()
			return
		}
		var wg sync.WaitGroup
		wg.Add(count)
		for i := 1; i <= count; i++ {
			go worker(&wg, i)
		}
		wg.Wait()
	} else {
		help()
	}
}

func commandRouter(s string, ) {
	if s == "single" {
		single()
	} else if s == "multiple" {
		multiple()
	} else {
		help()
	}
}

//show cli help
func help() {
	fmt.Println("    GOLANG Stress Test :: Arash Rasoulzadeh ")
	fmt.Println("------------------------------------------------")
	fmt.Println(" single            ::  run a single application")
	fmt.Println(" multiple {counts} ::  run with {counts} threads")
	fmt.Println("------------------------------------------------")
}

//main function
func main() {
	//worker()
	input := os.Args
	if len(input) >= 2 {
		clearLog()
		command := input[1]
		commandRouter(command)
	} else {
		help()
	}
}
