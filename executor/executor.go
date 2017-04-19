package executor

import (
	"bufio"
	"github.com/creichlin/pentaconta/declaration"
	"github.com/creichlin/pentaconta/logger"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

type Executor struct {
	service         *declaration.Service
	logs            *logger.Logger
	binary          string
	cmd             *exec.Cmd
	terminations    int
	terminationLock *sync.Mutex
}

func NewExecutor(service *declaration.Service, logs *logger.Logger) (*Executor, error) {

	binary, err := exec.LookPath(service.Executable)
	if err != nil {
		return nil, err
	}

	return &Executor{
		service:         service,
		binary:          binary,
		logs:            logs,
		terminationLock: &sync.Mutex{},
	}, nil
}

func (e *Executor) Start() {
	for {
		e.startService()
		// wait one second before restarting
		time.Sleep(time.Millisecond * 1000)
	}
}

func (e *Executor) Stop() {
	go func() {

		e.terminationLock.Lock()
		terminations := e.terminations
		e.cmd.Process.Signal(os.Interrupt)
		e.terminationLock.Unlock()

		for i := 0; i < 10; i++ {
			time.Sleep(time.Millisecond * 300)
			e.terminationLock.Lock()
			terminated := terminations != e.terminations
			e.terminationLock.Unlock()

			if terminated {
				e.logs.Log(logger.NewLog(logger.PENTACONTA, e.service.Name, "Sigint worked"))
				return
			}
		}

		e.logs.Log(logger.NewLog(logger.PENTACONTA, e.service.Name, "Sigint did not work, send kill"))
		e.cmd.Process.Signal(os.Kill)
	}()
}

func (e *Executor) startService() {
	e.cmd = exec.Command(e.binary, e.service.Arguments...)
	e.cmd.Env = os.Environ()

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
		log.Print(execErr)
	}

	blocker := make(chan int)
	go func() {
		// this will terminate itself when the command exits
		buffStdout := bufio.NewReader(stdout)
		for {
			line, err := buffStdout.ReadString('\n')
			if err == nil || line != "" {
				e.logs.Log(logger.NewLog(logger.STDOUT, e.service.Name, line))
			}
			if err != nil {
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
			line, err := buffStderr.ReadString('\n')
			if err == nil || line != "" {
				e.logs.Log(logger.NewLog(logger.STDERR, e.service.Name, line))
			}
			if err != nil {
				blocker <- 1
				// e.logs.Log(logger.NewLog(logger.PENTACOTA, e.service.Name, "Stderr ended"))
				return
			}
		}
	}()

	// we block till the stdout AND stderr reader are finished
	<-blocker
	<-blocker
	err = e.cmd.Wait()
	msg := "Terminated service"
	if err != nil {
		msg += " with " + err.Error()
	}
	e.logs.Log(logger.NewLog(logger.PENTACONTA, e.service.Name, msg))
	e.terminationLock.Lock()
	e.terminations += 1
	e.terminationLock.Unlock()
}
