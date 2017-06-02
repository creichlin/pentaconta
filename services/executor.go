package services

import (
	"bufio"
	"gitlab.com/creichlin/pentaconta/declaration"
	"gitlab.com/creichlin/pentaconta/logger"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"syscall"
	"time"
)

type Executor struct {
	service         *declaration.Service
	name            string
	logs            logger.Logger
	binary          string
	cmd             *exec.Cmd
	terminations    int
	terminationLock *sync.Mutex
	running         int
	runningLock     *sync.Mutex
}

func NewExecutor(name string, service *declaration.Service, logs logger.Logger) (*Executor, error) {

	binary, err := exec.LookPath(service.Executable)
	if err != nil {
		return nil, err
	}

	return &Executor{
		name:            name,
		service:         service,
		binary:          binary,
		logs:            logs,
		terminationLock: &sync.Mutex{},
		runningLock:     &sync.Mutex{},
	}, nil
}

func (e *Executor) Log(level int, message string) {
	e.logs.Log(logger.NewLog(level, e.name, e.terminations, message))
}

func (e *Executor) Start() {
	for {
		e.startService()
		// wait one second before restarting
		time.Sleep(time.Millisecond * 1000)
	}
}

func (e *Executor) IsRunning() bool {
	e.runningLock.Lock()
	isr := e.running > 0
	e.runningLock.Unlock()
	return isr
}

func (e *Executor) Stop() {
	if !e.IsRunning() {
		return
	}

	go func() {
		e.terminationLock.Lock()
		terminations := e.terminations
		var sig os.Signal

		if runtime.GOOS == "windows" {
			sig = os.Interrupt
		} else {
			sig = syscall.SIGABRT
		}
		e.cmd.Process.Signal(sig)
		e.terminationLock.Unlock()

		for i := 0; i < 10; i++ {
			time.Sleep(time.Millisecond * 300)
			e.terminationLock.Lock()
			terminated := terminations != e.terminations
			e.terminationLock.Unlock()

			if terminated {
				e.Log(logger.PENTACONTA, "Sigint worked")
				return
			}
		}

		e.Log(logger.PENTACONTA, "Sigint did not work, send kill")
		e.cmd.Process.Signal(os.Kill)
	}()
}

func (e *Executor) startService() {
	e.cmd = exec.Command(e.binary, e.service.Arguments...)
	e.cmd.Env = os.Environ()

	if e.service.WorkingDir != "" {
		e.cmd.Dir = e.service.WorkingDir
	}

	stdout, err := e.cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	stderr, err := e.cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	execErr := e.cmd.Start()
	if execErr != nil {
		e.Log(logger.PENTACONTA, execErr.Error())
		return
	}

	e.Log(logger.PENTACONTA, "Started service")

	blocker := make(chan int)
	go func() {
		// this will terminate itself when the command exits
		buffStdout := bufio.NewReader(stdout)
		for {
			line, errstdo := buffStdout.ReadString('\n')
			if errstdo == nil || line != "" {
				e.Log(logger.STDOUT, line)
			}
			if errstdo != nil {
				blocker <- 1
				// e.logs.Log(logger.NewLog(logger.PENTACOTA, e.service.Name, "Stdout ended"))
				return
			}
		}
	}()

	go func() {
		// this will terminate itself when the command exits
		buffStderr := bufio.NewReader(stderr)
		for {
			line, errstde := buffStderr.ReadString('\n')
			if errstde == nil || line != "" {
				e.Log(logger.STDERR, line)
			}
			if errstde != nil {
				blocker <- 1
				// e.logs.Log(logger.NewLog(logger.PENTACOTA, e.service.Name, "Stderr ended"))
				return
			}
		}
	}()

	e.runningLock.Lock()
	e.running++
	e.runningLock.Unlock()

	// we block till the stdout AND stderr reader are finished
	<-blocker
	<-blocker
	err = e.cmd.Wait()
	msg := "Terminated service"
	if err != nil {
		msg += " with " + err.Error()
	}
	e.Log(logger.PENTACONTA, msg)

	e.terminationLock.Lock()
	e.terminations += 1
	e.terminationLock.Unlock()

	e.runningLock.Lock()
	e.running--
	e.runningLock.Unlock()
}
