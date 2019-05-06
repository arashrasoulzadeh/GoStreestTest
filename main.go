package main

import (
	"bytes"
	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/arashrasoulzadeh/GoStreestTest/confutils"
	"github.com/schollz/progressbar"
)

type headersType map[string]string

const defaultPathToConfigFile = "config.yaml"

// send http call
func SendHttpRequest(method string, url string, body string, headers headersType) (*http.Response, error) {
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

		if len(headers) != 0 {
			for header := range headers {
				req.Header.Set(header, headers[header])
			}
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		return resp, err
	}

	fmt.Println("No valid method was provided, using default method GET:")
	resp, err := http.Get(url)
	return resp, err
}

//recover if error occured in MakeRequest
func MakeRequestRecovery(wg *sync.WaitGroup, bar *progressbar.ProgressBar) {
	defer wg.Done()
	bar.Add(1)
	if r := recover(); r != nil {
		fmt.Println("Recovered in f", r)
	}
}

func getError(err string) string {
	if strings.Contains(err, "refused") {
		return "REFUSED"
	}
	if strings.Contains(err, "reset") {
		return "RESET"
	}
	return "UNKNOWN"
}

//make request handler
func MakeRequest(thread int, url string, method string, body string, headers headersType, id int, wg *sync.WaitGroup, bar *progressbar.ProgressBar, f *os.File, ferr error, values *list.List) {
	defer MakeRequestRecovery(wg, bar)
	start := time.Now()
	resp, err := SendHttpRequest(method, url, body, headers)
	duration := time.Since(start).Seconds()
	if err != nil {
		// handle the error, often:
		msg := fmt.Sprintf("%d,%d,%d,%f,%s,%s\n",
			thread,
			id,
			-1,
			duration, "-1", fmt.Sprint(getError(err.Error())))
		values.PushFront(msg)
		return
	}
	msg := fmt.Sprintf("%d,%d,%d,%f,%s,%s\n",
		thread,
		id,
		resp.StatusCode,
		duration, "0", "NA")
	values.PushFront(msg)
}

//write to the log
func writeToLog(values *list.List, f *os.File) {
	//fmt.Println(values.Len())
	for temp := values.Front(); temp != nil; temp = temp.Next() {
		f.WriteString(fmt.Sprint(temp.Value))
	}

}

//clear log for use
func clearLog() {
	message := []byte("thread,id,code,duration,error,error_desc\n")
	err := ioutil.WriteFile("log.csv", message, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

//worker
func worker(mainWaitGroup *sync.WaitGroup, thread int, total int, bar *progressbar.ProgressBar, values *list.List, f *os.File, ferr error) {
	defer func() {
		defer mainWaitGroup.Done()
	}()
	var wg sync.WaitGroup
	var c confutils.Conf
	c.GetConf(defaultPathToConfigFile)
	_ = time.Now()
	wg.Add(c.Hits)
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
		go MakeRequest(thread, c.Route, c.Method, body, c.Headers, i, &wg, bar, f, ferr, values)
	}
	wg.Wait()
}

func single(values *list.List) {
	f, err := os.OpenFile("log.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer func() {
		f.Close()
	}()
	var c confutils.Conf
	c.GetConf(getConfig())
	bar := *progressbar.New(c.Hits * 1)
	bar.RenderBlank()
	var wg sync.WaitGroup
	wg.Add(1)
	worker(&wg, 1, 1, &bar, values, f, err)
	writeToLog(values, f)
	bar.Finish()
	fmt.Print("\n")
}

func multiple(values *list.List) {

	input := os.Args
	if len(input) >= 3 {
		f, err := os.OpenFile("log.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}
		defer func() {
			f.Close()
		}()
		counts := input[2]
		count, err := strconv.Atoi(counts)
		if err != nil {
			help()
			return
		}
		var c confutils.Conf
		c.GetConf(getConfig())
		bar := *progressbar.New(c.Hits * count)
		bar.RenderBlank()
		var wg sync.WaitGroup
		wg.Add(count)
		for i := 1; i <= count; i++ {
			go worker(&wg, i, count, &bar, values, f, err)
		}
		wg.Wait()
		writeToLog(values, f)
		bar.Finish()
		fmt.Print("\n")
	} else {

		help()
	}
}

func getConfig() string {
	input := os.Args
	if len(input) >= 2 {
		command := input[1]
		if command == "single" {
			if len(input) == 3 {
				return input[2]
			}
		}
		if command == "multiple" {
			if len(input) == 4 {
				return input[3]
			}
		}
	}
	return defaultPathToConfigFile
}

func commandRouter(s string, values *list.List) {
	if s == "single" {
		single(values)

	} else if s == "multiple" {
		multiple(values)
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
	values := list.New()
	input := os.Args
	if len(input) >= 2 {
		clearLog()
		command := input[1]
		commandRouter(command, values)
	} else {
		help()
	}
}
