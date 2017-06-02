package main

import (
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


func TestCustomSignals(t *testing.T) {
	sr := NewServiceRunner("testdata/integration/custom_signal.yaml")
	defer sr.Close()
	ioutil.WriteFile(filepath.Join(sr.Dir, "aborted.txt"), []byte("data"), 0644)
	go sr.Start(time.Millisecond * 1500)
	time.Sleep(time.Millisecond * 500)
	ioutil.WriteFile(filepath.Join(sr.Dir, "aborted.txt"), []byte("datax"), 0644)
	time.Sleep(time.Millisecond * 500)

}



func TestRuns(t *testing.T) {
	sr := NewServiceRunner("testdata/integration/1s.yaml")
	defer sr.Close()
	sr.Start(time.Second * 2)
	stats, err := ioutil.ReadFile(filepath.Join(sr.Dir, "stats.json"))
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


type ServiceRunner struct {
	data interface{}
	oldWorkingDir string
	Dir string
}

func (s *ServiceRunner)Start(duration time.Duration) {
	runWithDeclaration(s.data, duration)
}

func (s *ServiceRunner)Close() {
	os.RemoveAll(s.Dir)
	os.Chdir(s.oldWorkingDir)
}

func NewServiceRunner(configPath string) *ServiceRunner {
	dir, err := ioutil.TempDir("", "pentacontatest")
	if err != nil {
		log.Fatal(err)
	}

	cd, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	data := interface{}(nil)

	err = yaml.Unmarshal(cd, &data)
	if err != nil {
		panic(err)
	}

	oldWorkingDir, _ := os.Getwd()
	os.Chdir(dir)

	return &ServiceRunner{
		data: data,
		oldWorkingDir: oldWorkingDir,
		Dir: dir,
	}
}