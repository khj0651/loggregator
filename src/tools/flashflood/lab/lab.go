package lab

import (
	"sync"
	"tools/flashflood/experiment"
)

type Experiment interface {
	Start(stop <-chan struct{})
	Results() experiment.Results
}

type experimentFactory func() Experiment

type Lab struct {
	exp           Experiment
	stop chan struct{}
	m sync.Mutex
	newExperiment experimentFactory
}

func New(f experimentFactory) *Lab {
	return &Lab{
		newExperiment: f,
	}
}

func (l *Lab) Start() {
	l.m.Lock()
	defer l.m.Unlock()

	if l.stop != nil {
		close(l.stop)
	}

	l.exp = l.newExperiment()
	l.stop = make(chan struct{})
	go l.exp.Start(l.stop)
}

func (l* Lab) Stop() {
	l.m.Lock()
	defer l.m.Unlock()

	close(l.stop)
	l.stop = nil
}

func (l* Lab) Results() experiment.Results {
	l.m.Lock()
	defer l.m.Unlock()

	if (l.exp == nil) {
		return []experiment.Result{}
	}
	return l.exp.Results()
}