package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ghodss/yaml"
	"gitlab.com/creichlin/pentaconta/declaration"
	"gitlab.com/creichlin/pentaconta/evaluation"
	"gitlab.com/creichlin/pentaconta/logger"
	"gitlab.com/creichlin/pentaconta/services"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	configName, help := readConfigParams()
	if help {
		printHelp()
		return
	}
	location, err := probeLocation(configName)

	if err != nil {
		log.Fatal(err)
	}

	data, err := readData(location)
	if err != nil {
		log.Fatal(err)
	}

	runWithDeclaration(data, time.Hour*24*356*100)
}

func runWithDeclaration(data interface{}, duration time.Duration) {
	dec, err := declaration.Parse(data)
	if err != nil {
		log.Fatal(err)
	}

	var logs logger.Logger
	var eval *evaluation.Collector
	if dec.Stats != nil {
		eval = evaluation.EvaluationCollector()
		logs = logger.NewSplitLogger(
			logger.NewStdoutLogger(),
			eval,
		)
	} else {
		logs = logger.NewStdoutLogger()
	}

	services := &services.Services{
		Logs:        logs,
		Executors:   map[string]*services.Executor{},
		FSListeners: map[string]*services.FSListener{},
	}

	createAndStartExecutors(services, dec)
	createAndStartFsTriggers(services, dec)

	startTime := time.Now()

	for time.Now().Before(startTime.Add(duration)) {
		time.Sleep(time.Second * 1)

		if eval != nil {
			stats := eval.Status(dec.Stats.Seconds)
			bin, err := json.MarshalIndent(stats, "", "  ")
			if err != nil {
				panic(err)
			}
			err = ioutil.WriteFile(dec.Stats.File, bin, 0644)
			if err != nil {
				panic(err)
			}
		}
	}
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

func readConfigParams() (string, bool) {
	var configName string
	var help bool
	executable := filepath.Base(os.Args[0])
	flags := flag.NewFlagSet("pentacota", flag.ContinueOnError)
	flags.StringVar(&configName, "config", executable, "name of config file to use, no .yaml or .json extension.")
	flags.BoolVar(&help, "help", false, "Print help text and exit")
	flags.Parse(os.Args[1:])
	return configName, help
}

func printHelp() {
	fmt.Println("pentaconta help")
	fmt.Println(declaration.Doc())
}

func readData(file string) (interface{}, error) {
	binData, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	data := interface{}(nil)

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
		abspath := filepath.Join("/etc", path)
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
