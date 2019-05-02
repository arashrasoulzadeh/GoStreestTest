package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os/exec"
	"sync"
)

type scenarioConf struct {
	Label      string `yaml:"label"`
	ConfigFile string `yaml:"config"`
}

//strucure of the yaml file
type scenario struct {
	Scenarios []scenarioConf `yaml:"scenarios"`
}

const defaultPathToScenarioFile = "scenario.yaml"

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

func run(ConfigFile string,wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	app := "go"
	//app := "buah"

	arg0 := "run"
	arg1 := "main.go"
	arg2 := "single"

	cmd := exec.Command(app, arg0, arg1, arg2)
	stdout, err := cmd.Output()

	if err != nil {
		println(err.Error())
		return
	}

	print(string(stdout))

}
func main() {
	var wg sync.WaitGroup
	fmt.Println("GoStressTest Composer")
	var c scenario
	c.getScenario(defaultPathToScenarioFile)
	for _, element := range c.Scenarios {
		wg.Add(1);
		fmt.Println("\n\n running config file")
		run(element.ConfigFile,&wg)
	}
	defer func() {
		fmt.Println("halt.")
	}()

	wg.Wait()

}