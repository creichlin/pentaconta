package services

import (
	"fmt"
	"gitlab.com/creichlin/pentaconta/declaration"
	"gitlab.com/creichlin/pentaconta/logger"
	"reflect"
	"strings"
	"testing"
	"time"
)

func createPCExecutor(name string) (*[]string, *Executor, error) {
	logs := []string{}

	callback := func(t time.Time, level, service string, instance int, message string) {
		logs = append(logs, fmt.Sprintf("%v %v%v %v", level, service, instance, message))
	}

	service := &declaration.Service{
		Executable: "pc_" + name,
		Arguments:  []string{},
	}

	executor, err := NewExecutor("foo", service, logger.NewCallLogger(callback))
	if err != nil {
		return nil, nil, err
	}

	go executor.Start()

	return &logs, executor, nil
}

func createExecutor(t *testing.T, service *declaration.Service) (*[]string, *Executor) {
	logs := []string{}

	callback := func(t time.Time, level, service string, instance int, message string) {
		logs = append(logs, fmt.Sprintf("%v %v%v %v", level, service, instance, message))
	}

	executor, err := NewExecutor("foo", service, logger.NewCallLogger(callback))
	if err != nil {
		t.Error(err)
		return nil, nil
	}

	go executor.Start()

	return &logs, executor
}

func TestWorkingDir(t *testing.T) {
	if logs, _ := createExecutor(t, &declaration.Service{
		Executable: "pc_info",
		WorkingDir: "/etc",
	}); logs != nil {
		time.Sleep(time.Millisecond * 100)
		if !reflect.DeepEqual(*logs, []string{
			"PEN foo0 Started service",
			"OUT foo0 arguments: ",
			"OUT foo0 cwd: /etc",
		}) {
			t.Errorf("Wrong messages logged, %v", strings.Join(*logs, "\",\n\""))
		}
	}

}

func TestCrash(t *testing.T) {
	logs, _, err := createPCExecutor("unstable")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Millisecond * 1500)
	if !reflect.DeepEqual(*logs, []string{
		"PEN foo0 Started service",
		"OUT foo0 Unstable main started",
		"PEN foo0 Terminated service with exit status 2",
		"PEN foo1 Started service",
		"OUT foo1 Unstable main started",
		"PEN foo1 Terminated service with exit status 2",
	}) {
		t.Errorf("Wrong messages logged, %v", strings.Join(*logs, "\",\n\""))
	}
}

func TestLogs(t *testing.T) {
	logs, _, err := createPCExecutor("stable")
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
	logs, executor, err := createPCExecutor("stable")
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
