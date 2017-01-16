package store

import (
	"time"

	"github.com/satori/go.uuid"

	"github.com/InnovaCo/serve-runner/config"
	"github.com/InnovaCo/serve-runner/logger"
	"github.com/InnovaCo/serve-runner/sensors"
	"github.com/InnovaCo/serve-runner/tasks"
)

func init() {
	registry("default", Store(&defaultStore{}))
}

type defaultStore struct {
	flows []InstFlow
}

func (p *defaultStore) Init() error {
	p.flows = []InstFlow{}
	return nil
}

func (p *defaultStore) AddFlow(trigger sensors.Trigger, cfg config.FlowConfig) (InstFlow, error) {
	tp := []tasks.Task{}
	cp := []tasks.Context{}
	for _, t := range cfg.Tasks() {
		tp = append(tp, tasks.WorkerPool[t.Task()].Task)
		cp = append(cp, tasks.Context{Kwargs: t.Context()})
	}
	instFlow := InstFlow{
		Uuid:        uuid.NewV4().String(),
		Trigger:     trigger,
		Exchanger:   tasks.GetExchanger(),
		TaskPool:    tp,
		ContextPool: cp,
		Results:     []tasks.Result{},
		CurrentTask: TaskInit,
		lastModify:  time.Now()}
	p.flows = append(p.flows, instFlow)
	logger.Log.Debugf("Add flow: %v", instFlow.Uuid)
	logger.Log.Debugf("Current flows: %v", p.flows)

	return instFlow, nil
}

func (p *defaultStore) SaveFlow(flow InstFlow) error {
	for i, f := range p.flows {
		if flow.Uuid == f.Uuid {
			logger.Log.Debugf("%v save\n", f.Uuid)
			p.flows[i] = flow
			p.flows[i].lastModify = time.Now()
			return nil
		}
	}
	p.flows = append(p.flows, flow)
	return nil
}

func (p *defaultStore) GetFlows() []InstFlow {
	activeFlows := []InstFlow{}
	for _, flow := range p.flows {
		if len(flow.Results) != len(flow.TaskPool) {
			activeFlows = append(activeFlows, flow)
		}
	}
	return activeFlows
}
