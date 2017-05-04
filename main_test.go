package main

import (
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestReadDefaultConfigName(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()

	os.Args = []string{oldArgs[0]}
	name, _ := readConfigParams()
	if name != "pentaconta.test" {
		t.Errorf("Default config name should be name of executable but is %v", name)
	}
}

func TestReadCustomConfigName(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()

	os.Args = []string{oldArgs[0], "-config", "custom-name"}
	name, _ := readConfigParams()
	if name != "custom-name" {
		t.Errorf("Custom config name should be custom-name but is %v", name)
	}
}

func TestRuns(t *testing.T) {
	dir, err := ioutil.TempDir("", "pentacontatest")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
	defer os.RemoveAll(dir) // clean up

	cd, err := ioutil.ReadFile("testdata/integration/1s.yaml")
	if err != nil {
		panic(err)
	}

	data := interface{}(nil)

	err = yaml.Unmarshal(cd, &data)
	if err != nil {
		panic(err)
	}

	os.Chdir(dir)

	runWithDeclaration(data, time.Second*2)

	stats, err := ioutil.ReadFile(filepath.Join(dir, "stats.json"))
	if err != nil {
		panic(err)
	}

	expected := `{
  "Samples": 2,
  "Services": {
    "pc_stable": {
      "Errors": 0,
      "Crashes": 0,
      "Terminations": 1,
      "Logs": 2
    }
  }
}`

	if string(stats) != expected {
		t.Errorf("Stats is not as expected:\nCurrent\n%v\nexpected:\n%v", string(stats), expected)
	}
}
