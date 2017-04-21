package services

import (
	"fmt"
	"github.com/creichlin/pentaconta/declaration"
	"github.com/creichlin/pentaconta/logger"
	"reflect"
	"testing"
	"time"
)

func createPcStableExecutor() (*[]string, *Executor, error) {
	logs := []string{}

	callback := func(t time.Time, level, service, message string) {
		logs = append(logs, fmt.Sprintf("%v %v %v", level, service, message))
	}

	service := &declaration.Service{
		Name:       "foo",
		Executable: "pc_stable",
		Arguments:  []string{},
	}

	executor, err := NewExecutor(service, logger.NewCallLogger(callback))
	if err != nil {
		return nil, nil, err
	}

	go executor.Start()

	return &logs, executor, nil
}

func TestLogs(t *testing.T) {
	logs, _, err := createPcStableExecutor()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Millisecond * 1500)
	if !reflect.DeepEqual(*logs, []string{
		"OUT foo Stable main started",
		"OUT foo I'm doing fine",
	}) {
		t.Errorf("Wrong messages logged, %v", logs)
	}
}

func TestStop(t *testing.T) {
	logs, executor, err := createPcStableExecutor()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Millisecond * 500)
	executor.Stop()
	time.Sleep(time.Millisecond * 500)

	if !reflect.DeepEqual(*logs, []string{
		"OUT foo Stable main started",
		"PEN foo Terminated service with signal: interrupt",
		"PEN foo Sigint worked",
	}) {
		t.Errorf("Wrong messages logged, %v", logs)
	}
}
