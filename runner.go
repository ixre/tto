package tto

import "sync"

type RunnerCalc struct {
	TablesRuns bool
	states     map[string]bool
	_mux       *sync.RWMutex
}

func NewRunnerCalc() RunnerCalc {
	return RunnerCalc{
		TablesRuns: false,
		_mux:       &sync.RWMutex{},
		states:     make(map[string]bool, 0),
	}
}

func (c *RunnerCalc) SignState(key string, b bool) {
	c._mux.Lock()
	defer c._mux.Unlock()
	c.states[key] = b
}

func (c *RunnerCalc) State(key string) bool {
	c._mux.RLock()
	defer c._mux.RUnlock()
	return c.states[key]
}
