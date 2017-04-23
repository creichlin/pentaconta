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
		Executable: "pc_stable",
		Arguments:  []string{},
	}

	executor, err := NewExecutor("foo", service, logger.NewCallLogger(callback))
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
		"PEN foo0 Started service",
		"OUT foo0 Stable main started",
		"OUT foo0 I'm doing fine",
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
		"PEN foo0 Started service",
		"OUT foo0 Stable main started",
		"PEN foo0 Terminated service with signal: interrupt",
		"PEN foo1 Sigint worked",
	}) {
		t.Errorf("Wrong messages logged, %v", logs)
	}
}
