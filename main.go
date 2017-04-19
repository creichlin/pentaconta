package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/creichlin/pentaconta/declaration"
	"github.com/creichlin/pentaconta/executor"
	"github.com/creichlin/pentaconta/logger"
	"github.com/ghodss/yaml"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {

	location, err := probeLocation(readConfigName())

	if err != nil {
		log.Fatal(err)
	}

	dataMap, err := readData(location)
	if err != nil {
		log.Fatal(err)
	}

	data := &declaration.Root{}

	err = mapstructure.Decode(dataMap, data)
	if err != nil {
		log.Fatal(err)
	}

	logs := logger.NewLogger()

	executors := []*executor.Executor{}

	for _, service := range data.Services {

		executor, err := executor.NewExecutor(service, logs)
		if err != nil {
			panic(err)
		}

		executors = append(executors, executor)
	}

	for _, executor := range executors {
		go executor.Start()
	}

	time.Sleep(time.Second * 5)

	for _, executor := range executors {
		executor.Stop()
	}

	time.Sleep(time.Second * 5)

}

func readConfigName() string {
	var configName string
	executable := filepath.Base(os.Args[0])
	flags := flag.NewFlagSet("pentacota", flag.ContinueOnError)
	flags.StringVar(&configName, "config", executable, "name of config file to use, no .yaml or .json extension.")
	flags.Parse(os.Args[1:])
	return configName
}

func readData(file string) (map[string]interface{}, error) {
	binData, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{}

	if strings.HasSuffix(file, ".json") {
		err = json.Unmarshal(binData, &data)
		if err != nil {
			return nil, err
		}
		return data, nil
	} else if strings.HasSuffix(file, ".yaml") {
		err = yaml.Unmarshal(binData, &data)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	panic("Returned path must have .json or .yaml extension")
}

func probeLocation(path string) (string, error) {
	locations := []string{}
	if filepath.IsAbs(path) {
		locations = append(locations, path+".json")
		locations = append(locations, path+".yaml")
	} else {
		wd, err := os.Getwd()
		if err == nil {
			abspath := filepath.Join(wd, path)
			locations = append(locations, abspath+".json")
			locations = append(locations, abspath+".yaml")
		}
		abspath := filepath.Join("/etc", "argo", path)
		locations = append(locations, abspath+".json")
		locations = append(locations, abspath+".yaml")
	}

	for _, location := range locations {
		_, err := ioutil.ReadFile(location)
		if err == nil {
			return location, nil
		}
	}

	return "", fmt.Errorf("Could not find config file in locations %v", locations)
}
