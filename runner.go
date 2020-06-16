package tto

type RunnerCalc struct {
	TablesRuns       bool
	TablesStates     map[string]bool
	GroupTablesState map[string]bool
}

func NewRunnerCalc()RunnerCalc{
	return RunnerCalc{
		TablesRuns:   false,
		TablesStates: make(map[string]bool,0),
		GroupTablesState:make(map[string]bool,0),
	}
}