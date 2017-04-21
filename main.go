package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/creichlin/pentaconta/declaration"
	"github.com/creichlin/pentaconta/logger"
	"github.com/creichlin/pentaconta/services"
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

	services := &services.Services{
		Logs:        logger.NewStdoutLogger(),
		Executors:   map[string]*services.Executor{},
		FSListeners: map[string]*services.FSListener{},
	}

	createAndStartExecutors(services, data)
	createAndStartFsTriggers(services, data)

	time.Sleep(time.Second * 20)

}

func createAndStartFsTriggers(svs *services.Services, data *declaration.Root) {
	for name, fsTrigger := range data.FSTriggers {
		fsListener, err := services.NewFSListener(name, fsTrigger, svs)
		if err != nil {
			panic(err)
		}
		svs.FSListeners[name] = fsListener
		go func() {
			err := fsListener.Start()
			log.Fatal(err)
		}()
	}
}

func createAndStartExecutors(svs *services.Services, data *declaration.Root) {
	for name, service := range data.Services {
		executor, err := services.NewExecutor(name, service, svs.Logs)
		if err != nil {
			panic(err)
		}

		svs.Executors[name] = executor
		go executor.Start()
	}
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
