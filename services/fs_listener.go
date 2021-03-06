package services

import (
	"github.com/creichlin/gutil/run"
	"gitlab.com/creichlin/pentaconta/declaration"
)

type FSListener struct {
	name     string
	trigger  *declaration.FSTrigger
	services *Services
	ftr      *run.FileTriggerRunner
}

func NewFSListener(name string, trigger *declaration.FSTrigger, services *Services) (*FSListener, error) {
	listener := &FSListener{
		name:     name,
		trigger:  trigger,
		services: services,
	}
	listener.ftr = run.NewFileTriggerRunner(trigger.Path, false, listener.changed)

	return listener, nil
}

func (f *FSListener) changed(event, path string) error {
	for _, service := range f.trigger.Services {
		if f.trigger.Signal == "" {
			f.services.Executors[service].Stop()
		}else{
			f.services.Executors[service].Signal(f.trigger.Signal)
		}
	}
	return nil
}

func (f *FSListener) Start() error {
	return f.ftr.Start()
}
