package main

import (
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

type scenarioConf struct {
	Label      string `yaml:"label"`
	ConfigFile string `yaml:"config"`
	Sleep      string `yaml:"sleep"`
}

//strucure of the yaml file
type scenario struct {
	Scenarios []scenarioConf `yaml:"scenarios"`
}

const defaultPathToScenarioFile = "scenario.yaml"
const TYPE_ERROR = 0;
const TYPE_SUCCESS = 1;
const TYPE_INFO = 3;

//get the scenario
// dest: file destination / the path to scenario file
func (c *scenario) getScenario(dest string) *scenario {
	yamlFile, err := ioutil.ReadFile(dest)
	if err != nil {
		panic(fmt.Sprint("yamlFile.Get err   #%v ", err))
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func run(ConfigFile string, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	app := "go"
	arg0 := "run"
	arg1 := "main.go"
	arg2 := "single"
	arg3 := ConfigFile

	cmd := exec.Command(app, arg0, arg1, arg2, arg3)
	_, err := cmd.Output()

	if err != nil {
		coloredPrint("COMPOSER:ERROR  ", err.Error(), TYPE_ERROR)
		return
	}else{
		coloredPrint("COMPOSER:SUCCESS", ConfigFile, TYPE_SUCCESS)
	}

	//print(string(stdout))

}
func coloredPrint(colored string, msg string, msgtype int ) {
	label := color.New(color.FgGreen).SprintFunc()
	if msgtype==TYPE_ERROR {
		label = color.New(color.FgRed).SprintFunc()
	}
	if msgtype==TYPE_INFO {
		label = color.New(color.FgBlue).SprintFunc()
	}

	message := color.New(color.FgBlack).SprintFunc()
	fmt.Printf("%s :: %s\n", label(colored), message(msg))
}
func main() {
	var wg sync.WaitGroup
	fmt.Println("GoStressTest Composer")
	var c scenario
	c.getScenario(defaultPathToScenarioFile)
	for _, element := range c.Scenarios {
		wg.Add(1);
		coloredPrint("COMPOSER:RUN    ", element.Label, TYPE_INFO)
		run(element.ConfigFile, &wg)
		coloredPrint("COMPOSER:SLEEP  ", element.Sleep, TYPE_INFO)

		sleep, _ := strconv.Atoi(element.Sleep);
		time.Sleep(time.Duration(sleep) * time.Second)
	}
	defer func() {
		fmt.Println("halt.")
	}()

	wg.Wait()

}
