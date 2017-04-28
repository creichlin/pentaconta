package main

import (
	"os"
	"testing"
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
